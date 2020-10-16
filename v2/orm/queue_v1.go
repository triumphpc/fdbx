package orm

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/shestakovda/fdbx/v2"
	"github.com/shestakovda/fdbx/v2/db"
	"github.com/shestakovda/fdbx/v2/mvcc"
)

func Queue(id byte, cl Collection, opts ...Option) TaskCollection {
	q := v1Queue{
		id:      id,
		cl:      cl,
		options: newOptions(),
	}

	for i := range opts {
		opts[i](&q.options)
	}

	return &q
}

type v1Queue struct {
	options
	id byte
	cl Collection
}

func (q v1Queue) ID() byte { return q.id }

func (q v1Queue) Ack(tx mvcc.Tx, ids ...fdbx.Key) (err error) {
	keys := make([]fdbx.Key, 2*len(ids))

	for i := range ids {
		// Удаляем из неподтвержденных
		keys[2*i] = q.usrKey(ids[i].RPart(qWork))
		// Удаляем из индекса статусов задач
		keys[2*i+1] = q.usrKey(ids[i].RPart(iStat))
	}

	if err = tx.Delete(keys); err != nil {
		return ErrAck.WithReason(err)
	}

	tx.OnCommit(func(w db.Writer) error {
		// Уменьшаем счетчик задач в работе
		w.Increment(q.usrKey(qTotalWorkKey), int64(-len(ids)))
		return nil
	})
	return nil
}

func (q v1Queue) Pub(tx mvcc.Tx, when time.Time, ids ...fdbx.Key) (err error) {

	if when.IsZero() {
		when = time.Now()
	}

	delay := make([]byte, 8)
	binary.BigEndian.PutUint64(delay, uint64(when.UTC().UnixNano()))

	// Структура ключа:
	// db nsUser cl.id q.id qList delay uid = taskID
	pairs := make([]fdbx.Pair, 2*len(ids))
	for i := range ids {
		// Основная запись таски
		pairs[2*i] = fdbx.NewPair(
			// Случайная айдишка таски, чтобы не было конфликтов при одинаковом времени
			q.usrKey(ids[i].LPart(delay...).LPart(qList)),
			// Айдишку элемента очереди записываем в значение, именно она и является таской
			[]byte(ids[i]),
		)

		// Служебная запись в индекс состояний
		pairs[2*i+1] = fdbx.NewPair(
			q.usrKey(ids[i].LPart(iStat)),
			[]byte{StatusPublished},
		)
	}

	if err = tx.Upsert(pairs); err != nil {
		return ErrPub.WithReason(err)
	}

	// Особая магия - инкремент счетчика очереди, чтобы затриггерить подписчиков
	// А также инкремент счетчиков статистики очереди
	tx.OnCommit(func(w db.Writer) error {
		// Увеличиваем счетчик задач в ожидании
		w.Increment(q.usrKey(qTotalWaitKey), int64(len(ids)))
		// Триггерим обработчики забрать новые задачи
		w.Increment(q.usrKey(qTriggerKey), 1)
		return nil
	})

	return nil
}

func (q v1Queue) Sub(ctx context.Context, cn db.Connection, pack int) (<-chan fdbx.Pair, <-chan error) {
	res := make(chan fdbx.Pair)
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		defer close(res)
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					errc <- ErrSub.WithReason(err)
				} else {
					errc <- ErrSub.WithReason(fmt.Errorf("%+v", rec))
				}
			}
		}()

		hdlr := func() (err error) {
			var list []fdbx.Pair

			if list, err = q.SubList(ctx, cn, pack); err != nil {
				return
			}

			if len(list) == 0 {
				return
			}

			for i := range list {
				select {
				case res <- list[i]:
				case <-ctx.Done():
					return ErrSub.WithReason(ctx.Err())
				}
			}

			return nil
		}

		for {
			if err := hdlr(); err != nil {
				errc <- err
				return
			}
		}
	}()

	return res, errc
}

func (q v1Queue) SubList(ctx context.Context, cn db.Connection, pack int) (list []fdbx.Pair, err error) {
	if pack == 0 {
		return nil, nil
	}

	var pairs []fdbx.Pair
	var waiter fdbx.Waiter

	from := q.usrKey(fdbx.Key{qList})
	hdlr := func() (exp error) {
		var tx mvcc.Tx

		now := make([]byte, 8)
		binary.BigEndian.PutUint64(now, uint64(time.Now().UTC().UnixNano()))
		to := q.usrKey(fdbx.Key(now).LPart(qList))

		if tx, err = mvcc.Begin(cn); err != nil {
			return ErrSub.WithReason(err)
		}
		defer tx.Cancel()

		// Критически важно делать это в одной физической транзакции
		// Иначе остается шанс, что одну и ту же задачу возьмут в обработку два воркера
		return tx.Conn().Write(func(w db.Writer) (e error) {
			if pairs, e = tx.SeqScan(
				from, to, mvcc.Limit(pack),
				mvcc.Exclusive(q.onTaskWork),
				mvcc.Writer(w),
			); e != nil {
				return
			}

			if len(pairs) == 0 {
				// В этом случае не коммитим, т.к. по сути ничего не изменилось
				waiter = w.Watch(q.usrKey(qTriggerKey))
				return nil
			}

			ids := make([]fdbx.Key, len(pairs))

			for i := range pairs {
				if ids[i], e = pairs[i].WrapKey(q.waitKeyWrapper).Key(); e != nil {
					return
				}
			}

			if list, e = q.cl.Select(tx).PossibleByID(ids...).All(); e != nil {
				return
			}

			// Уменьшаем счетчик задач в ожидании
			w.Increment(q.usrKey(qTotalWaitKey), int64(-len(ids)))

			// Увеличиваем счетчик задач в ожидании
			w.Increment(q.usrKey(qTotalWorkKey), int64(len(ids)))

			// Логический коммит в той же физической транзакции
			// Это самый важный момент - именно благодаря этому перемещенные в процессе чтения
			// элементы очереди будут видны как перемещенные для других логических транзакций
			return tx.Commit(mvcc.Writer(w))
		})
	}

	for {
		if err = q.waitTask(ctx, waiter); err != nil {
			return
		}

		if err = hdlr(); err != nil {
			return nil, ErrSub.WithReason(err)
		}

		if len(list) > 0 {
			return list, nil
		}
	}
}

func (q v1Queue) Lost(tx mvcc.Tx, pack int) (list []fdbx.Pair, err error) {
	if pack == 0 {
		return nil, nil
	}

	var id []byte
	var pairs []fdbx.Pair

	key := q.usrKey(fdbx.Key{qWork})

	// Значения в этих парах - айдишки элементов коллекции
	if pairs, err = tx.SeqScan(key, key, mvcc.Limit(pack)); err != nil {
		return nil, ErrLost.WithReason(err)
	}

	if len(pairs) == 0 {
		return nil, nil
	}

	ids := make([]fdbx.Key, len(pairs))

	for i := range pairs {
		if id, err = pairs[i].Value(); err != nil {
			return nil, ErrLost.WithReason(err)
		}
		ids[i] = fdbx.Key(id)
	}

	if list, err = q.cl.Select(tx).PossibleByID(ids...).All(); err != nil {
		return nil, ErrLost.WithReason(err)
	}

	return list, nil
}

func (q v1Queue) Status(tx mvcc.Tx, ids ...fdbx.Key) (res map[string]byte, err error) {
	var val []byte
	var pair fdbx.Pair

	res = make(map[string]byte, len(ids))

	for i := range ids {
		status := StatusConfirmed

		if pair, err = tx.Select(q.usrKey(ids[i].RPart(iStat))); err == nil {
			if val, err = pair.Value(); err != nil {
				return nil, ErrStatus.WithReason(err)
			}
			if len(val) > 0 {
				status = val[0]
			}
		}

		res[ids[i].String()] = status
	}

	return res, nil
}

func (q v1Queue) Stat(tx mvcc.Tx) (wait, work int64, err error) {
	if err = tx.Conn().Read(func(r db.Reader) (exp error) {
		var val []byte

		if val, exp = r.Data(q.usrKey(qTotalWaitKey)).Value(); exp != nil {
			return
		}
		wait = int64(binary.LittleEndian.Uint64(val))

		if val, exp = r.Data(q.usrKey(qTotalWorkKey)).Value(); exp != nil {
			return
		}
		work = int64(binary.LittleEndian.Uint64(val))

		return nil
	}); err != nil {
		return 0, 0, ErrStat.WithReason(err)
	}

	return wait, work, nil
}

func (q v1Queue) waitTask(ctx context.Context, waiter fdbx.Waiter) (err error) {
	if waiter == nil {
		return nil
	}

	if ctx.Err() != nil {
		return ErrSub.WithReason(ctx.Err())
	}

	// Даже если waiter установлен, то при отсутствии других публикаций мы тут зависнем навечно.
	// А задачи, время которых настало, будут просрочены. Для этого нужен особый механизм обработки по таймауту.
	wctx, cancel := context.WithTimeout(ctx, q.options.punch)
	defer cancel()

	if err = waiter.Resolve(wctx); err != nil {
		return ErrSub.WithReason(err)
	}

	// Если запущено много обработчиков, все они рванут забирать события одновременно.
	// Чтобы избежать массовых конфликтов транзакций и улучшить распределение задач делаем небольшую
	// случайную задержку, в пределах 20 мс. Немного для человека, значительно для уменьшения конфликтов
	time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
	return nil
}

func (q v1Queue) usrKeyWrapper(key fdbx.Key) (fdbx.Key, error) { return q.usrKey(key), nil }

func (q v1Queue) usrKey(key fdbx.Key) fdbx.Key {
	clid := q.cl.ID()
	return key.LPart(byte(clid>>8), byte(clid), q.id)
}

func (q v1Queue) waitKeyWrapper(key fdbx.Key) (fdbx.Key, error) {
	return key.LSkip(12), nil
}

func (q v1Queue) wrkKeyWrapper(key fdbx.Key) (fdbx.Key, error) {
	return q.usrKey(key.LSkip(12).RPart(qWork)), nil
}

func (q v1Queue) onTaskWork(tx mvcc.Tx, p fdbx.Pair, w db.Writer) (exp error) {
	var key fdbx.Key

	if key, exp = p.Key(); exp != nil {
		return ErrSub.WithReason(exp)
	}

	// Удаление по ключу из основной очереди
	if exp = tx.Delete([]fdbx.Key{key}, mvcc.Writer(w)); exp != nil {
		return ErrSub.WithReason(exp)
	}

	pairs := []fdbx.Pair{
		// Вставка в коллекцию задач "в работе"
		p.WrapKey(q.wrkKeyWrapper),
		// Вставка в индекс статусов задач
		fdbx.NewPair(q.usrKey(key.LSkip(12).RPart(iStat)), []byte{StatusUnconfirmed}),
	}

	if exp = tx.Upsert(pairs, mvcc.Writer(w)); exp != nil {
		return ErrSub.WithReason(exp)
	}

	return nil
}