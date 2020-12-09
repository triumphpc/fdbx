package orm

import (
	"context"
	"time"

	"github.com/shestakovda/errx"
	"github.com/shestakovda/fdbx/v2"
	"github.com/shestakovda/fdbx/v2/db"
	"github.com/shestakovda/fdbx/v2/models"
	"github.com/shestakovda/fdbx/v2/mvcc"
)

func NewTable(id uint16, args ...Option) Table {
	return &v1Table{
		id:      id,
		options: getOpts(args),
	}
}

type v1Table struct {
	options
	id uint16
}

func (t *v1Table) ID() uint16 { return t.id }

func (t *v1Table) Select(tx mvcc.Tx) Query { return NewQuery(t, tx) }

func (t *v1Table) Cursor(tx mvcc.Tx, id string) (Query, error) { return loadQuery(t, tx, id) }

func (t *v1Table) Insert(tx mvcc.Tx, pairs ...fdbx.Pair) (err error) {
	return t.upsert(tx, true, pairs...)
}

func (t *v1Table) Upsert(tx mvcc.Tx, pairs ...fdbx.Pair) (err error) {
	return t.upsert(tx, false, pairs...)
}

func (t *v1Table) Delete(tx mvcc.Tx, keys ...fdbx.Key) (err error) {
	if len(keys) == 0 {
		return nil
	}

	cp := make([]fdbx.Key, len(keys))
	for i := range keys {
		cp[i] = WrapTableKey(t.id, keys[i])
	}

	if err = tx.Delete(cp, mvcc.OnDelete(t.onDelete)); err != nil {
		return ErrDelete.WithReason(err)
	}

	return nil
}

func (t *v1Table) upsert(tx mvcc.Tx, ins bool, pairs ...fdbx.Pair) (err error) {
	if len(pairs) == 0 {
		return nil
	}

	cp := make([]fdbx.Pair, len(pairs))
	for i := range pairs {
		if cp[i], err = newSysPair(tx, t.id, pairs[i]); err != nil {
			return ErrUpsert.WithReason(err)
		}
	}

	opts := []mvcc.Option{
		mvcc.OnUpdate(t.onUpdate),
		mvcc.OnDelete(t.onDelete),
	}

	if ins {
		opts = append(opts, mvcc.OnInsert(t.onInsert))
	}

	if err = tx.Upsert(cp, opts...); err != nil {
		return ErrUpsert.WithReason(err)
	}

	return nil
}

func (t *v1Table) onInsert(tx mvcc.Tx, pair fdbx.Pair) (err error) {
	if len(pair.Value()) > 0 {
		return ErrDuplicate.WithDebug(errx.Debug{
			"key": pair.Key().String(),
		})
	}

	return nil
}

func (t *v1Table) onUpdate(tx mvcc.Tx, pair fdbx.Pair) (err error) {

	if len(t.options.batchidx) == 0 {
		return nil
	}

	pval := pair.Value()
	pkey := pair.Key().Bytes()
	rows := make([]fdbx.Pair, 0, 32)
	var dict map[uint16][]fdbx.Key

	for k := range t.options.batchidx {
		if dict, err = t.options.batchidx[k](pval); err != nil {
			return
		}

		if len(dict) == 0 {
			continue
		}

		for idx, keys := range dict {
			for i := range keys {
				if keys[i] == nil || len(keys[i].Bytes()) == 0 {
					continue
				}
				rows = append(rows, fdbx.NewPair(WrapIndexKey(t.id, idx, keys[i]).RPart(pkey...), pkey))
			}
		}
	}

	if err = tx.Upsert(rows); err != nil {
		return ErrIdxUpsert.WithReason(err)
	}

	return nil
}

func (t *v1Table) onDelete(tx mvcc.Tx, pair fdbx.Pair) (err error) {
	var usr fdbx.Pair

	if len(t.options.batchidx) == 0 {
		return nil
	}

	// Здесь нам придется обернуть еще значение, которое возвращается, потому что оно не обработано уровнем ниже
	if usr, err = newUsrPair(tx, t.id, pair); err != nil {
		return ErrIdxDelete.WithReason(err)
	}

	uval := usr.Value()
	rows := make([]fdbx.Key, 0, 32)
	pkey := UnwrapTableKey(pair.Key()).Bytes()
	var dict map[uint16][]fdbx.Key

	for k := range t.options.batchidx {
		if dict, err = t.options.batchidx[k](uval); err != nil {
			return
		}

		if len(dict) == 0 {
			continue
		}

		for idx, keys := range dict {
			for i := range keys {
				if keys[i] == nil || len(keys[i].Bytes()) == 0 {
					continue
				}
				rows = append(rows, WrapIndexKey(t.id, idx, keys[i]).RPart(pkey...))
			}
		}
	}

	if err = tx.Delete(rows); err != nil {
		return ErrIdxDelete.WithReason(err)
	}

	return nil
}

func (t *v1Table) Autovacuum(ctx context.Context, cn db.Connection, args ...Option) {
	var err error

	opts := getOpts(args)
	tick := time.NewTicker(opts.vwait)
	defer tick.Stop()

	defer func() {
		// Перезапуск только в случае ошибки
		if err != nil {
			time.Sleep(time.Second)

			// И только если мы вообще можем еще запускать
			if ctx.Err() == nil {
				// Тогда стартуем заново и в s.wait ничего не ставим
				go t.Autovacuum(ctx, cn, args...)
				return
			}
		}
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

	for ctx.Err() == nil {

		if err = t.vacuumStep(cn); err != nil {
			return
		}

		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
	}
}

func (t *v1Table) vacuumStep(cn db.Connection) (err error) {
	var tx mvcc.Tx

	if tx, err = mvcc.Begin(cn); err != nil {
		return ErrVacuum.WithReason(err)
	}
	defer tx.Cancel()

	// Этот запрос очищает только данные. Для них должен быть обработчик очистки BLOB
	if err = tx.Vacuum(WrapTableKey(t.id, nil), mvcc.OnVacuum(t.onVacuum)); err != nil {
		return ErrVacuum.WithReason(err)
	}

	// Отдельно очистка всех индексов
	if err = tx.Vacuum(WrapIndexKey(t.id, 0, nil).RSkip(2)); err != nil {
		return ErrVacuum.WithReason(err)
	}

	// Отдельно очистка всех очередей
	if err = tx.Vacuum(WrapQueueKey(t.id, 0, nil, 0, nil).RSkip(3)); err != nil {
		return ErrVacuum.WithReason(err)
	}

	// Отдельно очистка всех блобов
	if err = tx.Vacuum(WrapBlobKey(t.id, nil)); err != nil {
		return ErrVacuum.WithReason(err)
	}

	// Отдельно очистка всех курсоров
	if err = tx.Vacuum(WrapQueryKey(t.id, nil)); err != nil {
		return ErrVacuum.WithReason(err)
	}

	return nil
}

func (t *v1Table) onVacuum(tx mvcc.Tx, p fdbx.Pair, w db.Writer) (err error) {
	var mod models.ValueT

	val := p.Value()

	if len(val) == 0 {
		return nil
	}

	models.GetRootAsValue(val, 0).UnPackTo(&mod)

	// Если значение лежит в BLOB, надо удалить
	if mod.Blob {
		if err = tx.DropBLOB(WrapBlobKey(t.id, fdbx.Bytes2Key(mod.Data)), mvcc.Writer(w)); err != nil {
			return ErrVacuum.WithReason(err)
		}
	}

	return nil
}
