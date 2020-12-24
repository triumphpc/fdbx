package mvcc

import (
	"github.com/shestakovda/fdbx/v2"
	"github.com/shestakovda/fdbx/v2/db"
)

func getOpts(args []Option) (o options) {
	o.packSize = 10

	for i := range args {
		args[i](&o)
	}

	return
}

type options struct {
	lock     bool
	reverse  bool
	limit    int
	packSize int
	from     fdbx.Key
	last     fdbx.Key
	onInsert PairHandler
	onUpdate PairHandler
	onDelete PairHandler
	onVacuum RowHandler
	onLock   RowHandler
	writer   db.Writer
}

func Lock() Option                    { return func(o *options) { o.lock = true } }
func Last(k fdbx.Key) Option          { return func(o *options) { o.last = k } }
func From(k fdbx.Key) Option          { return func(o *options) { o.from = k } }
func Limit(l int) Option              { return func(o *options) { o.limit = l } }
func Writer(w db.Writer) Option       { return func(o *options) { o.writer = w } }
func Reverse() Option                 { return func(o *options) { o.reverse = true } }
func OnInsert(hdl PairHandler) Option { return func(o *options) { o.onInsert = hdl } }
func OnUpdate(hdl PairHandler) Option { return func(o *options) { o.onUpdate = hdl } }
func OnDelete(hdl PairHandler) Option { return func(o *options) { o.onDelete = hdl } }
func OnVacuum(hdl RowHandler) Option  { return func(o *options) { o.onVacuum = hdl } }
func Exclusive(hdl RowHandler) Option { return func(o *options) { o.lock = true; o.onLock = hdl } }
func VacuumPack(size int) Option      { return func(o *options) { o.packSize = size } }
