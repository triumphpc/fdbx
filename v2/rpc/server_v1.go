package rpc

import (
	"context"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/shestakovda/errx"
	"github.com/shestakovda/fdbx/v2"
	"github.com/shestakovda/fdbx/v2/db"
	"github.com/shestakovda/fdbx/v2/mvcc"
	"github.com/shestakovda/fdbx/v2/orm"
)

func newServerV1(id uint16) Server {
	s := v1Server{
		data: orm.NewTable(id),
		list: make([]*endpoint, 0, 16),
	}

	return &s
}

type v1Server struct {
	data orm.Table
	list []*endpoint
	wait *sync.WaitGroup
	exit context.CancelFunc
}

func (s *v1Server) Endpoint(id uint16, hdl TaskHandler, args ...Option) (err error) {
	end := newEndpoint(id, s.data, hdl, args)

	if err = end.check(); err != nil {
		return
	}

	s.list = append(s.list, end)
	return nil
}

func (s *v1Server) Run(ctx context.Context, cn db.Connection, args ...Option) {
	var wctx context.Context

	wctx, s.exit = context.WithCancel(ctx)
	s.wait = new(sync.WaitGroup)
	s.wait.Add(len(s.list) + 1)

	for i := range s.list {
		go s.listen(wctx, cn, s.list[i])
	}

	go s.autovacuum(wctx, cn, args)
}

func (s *v1Server) Stop() {
	if s.exit != nil {
		s.exit()
	}

	s.wait.Wait()
}

func (s *v1Server) listen(ctx context.Context, cn db.Connection, end *endpoint) {
	var err error

	defer func() {
		// Перезапуск только в случае ошибки
		if err != nil {
			// Которая обработана и требует перезапуска
			if repeat, wait := end.OnListen(err); repeat {
				// Возможно, не сразу готовы обрабатывать снова
				if wait > 0 {
					time.Sleep(wait)
				}

				// И только если мы вообще можем еще запускать
				if ctx.Err() == nil {
					// Тогда стартуем заново и в s.wait ничего не ставим
					go s.listen(ctx, cn, end)
					return
				}
			}
		}

		// В остальных случаях, нечего ловить, закрываем ожидание
		s.wait.Done()
	}()

	// Отлавливаем панику и превращаем в ошибку
	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				err = ErrListen.WithReason(e)
			} else {
				err = ErrListen.WithDebug(errx.Debug{"panic": rec})
			}
		}
	}()

	// Собственный контекст для гарантированного завершения подписки в случае провала
	wctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pairs, errs := end.Queue.Sub(wctx, cn, 1)

	for pair := range pairs {
		// В случае ошибки при обработке задачи
		if res, exp := end.OnTask(pair); exp != nil {
			// Обрабатываем ошибку и если нужно, повторяем задачу
			if repeat, wait := end.OnError(exp); repeat {
				// Если не смогли повторить - это фиаско
				if err = end.repeat(cn, pair, wait); err != nil {
					err = ErrListen.WithReason(err)
					return
				}
			}
		} else {
			if err = end.confirm(cn, s.data, pair, res); err != nil {
				err = ErrListen.WithReason(err)
				return
			}
		}
	}

	for exp := range errs {
		if exp != nil && !errx.Is(exp, context.Canceled, context.DeadlineExceeded) {
			err = ErrListen.WithReason(exp)
			return
		}
	}

}

func (s *v1Server) autovacuum(ctx context.Context, cn db.Connection, args []Option) {
	var err error

	defer func() {
		// Перезапуск только в случае ошибки
		if err != nil {
			glog.Errorf("%+v", err)
			time.Sleep(time.Second)

			// И только если мы вообще можем еще запускать
			if ctx.Err() == nil {
				// Тогда стартуем заново и в s.wait ничего не ставим
				go s.autovacuum(ctx, cn, args)
				return
			}
		}

		// В остальных случаях, нечего ловить, закрываем ожидание
		s.wait.Done()
	}()

	// Отлавливаем панику и превращаем в ошибку
	defer func() {
		if rec := recover(); rec != nil {
			if e, ok := rec.(error); ok {
				err = ErrVacuum.WithReason(e)
			} else {
				err = ErrVacuum.WithDebug(errx.Debug{"panic": rec})
			}
		}
	}()

	opts := getOpts(args)
	wkey := mvcc.NewTxKeyManager().Wrap(orm.NewWatchKeyManager(s.data.ID()).Wrap(nil))
	tick := time.NewTicker(opts.vwait)
	defer tick.Stop()

	for ctx.Err() == nil {

		if err = cn.Write(func(w db.Writer) (exp error) {
			var val []byte
			var key fdbx.Key
			var when time.Time

			pairs := w.List(wkey, wkey, opts.vpack, false)()

			for i := range pairs {
				if val, exp = pairs[i].Value(); exp != nil {
					return
				}

				if when, exp = fdbx.Byte2Time(val); exp != nil {
					return
				}
				glog.Errorf("%s", when)

				if time.Since(when) < 24*time.Hour {
					continue
				}

				if key, exp = pairs[i].Key(); exp != nil {
					return
				}

				glog.Errorf("%s - %s", key, when)
				w.Delete(key)
			}

			return nil
		}); err != nil {
			err = ErrVacuum.WithReason(err)
			return
		}

		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
	}
}