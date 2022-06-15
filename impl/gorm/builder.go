package gorm

import (
	"context"
	"database/sql"
	"reflect"

	"gorm.io/gorm"

	"github.com/thoohv5/orm/standard"
)

// Add add
func (g *gDb) Add(ctx context.Context, value interface{}) error {
	return g.gDB().WithContext(ctx).Create(value).Error
}

// Update update
func (g *gDb) Update(ctx context.Context, attrs interface{}) error {
	return g.gDB().WithContext(ctx).Updates(attrs).Error
}

func (g *gDb) Session(sess *gorm.Session) standard.IBuilder {
	return CopyGDb(g.gDB().Session(sess), g.transmit()...)
}

// Delete delete
func (g *gDb) Delete(ctx context.Context, value interface{}, where ...interface{}) error {
	return g.gDB().WithContext(ctx).Delete(value, where...).Error
}

// Find find
func (g *gDb) Find(ctx context.Context, out interface{}, where ...interface{}) error {
	return g.TransformNoRecordErr(g.gDB().WithContext(ctx).First(out, where...).Error)
}

// Get get
func (g *gDb) Get(ctx context.Context, out interface{}, where ...interface{}) error {
	return g.gDB().WithContext(ctx).Find(out, where...).Error
}

func (g *gDb) GetGDB() *gorm.DB {
	return g.gDB()
}

// Count count
func (g *gDb) Count(ctx context.Context, value *int64) error {
	return g.TransformNoRecordErr(g.gDB().WithContext(ctx).Count(value).Error)
}

// Set set
func (g *gDb) Set(name string, value interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Set(name, value), g.transmit()...)
}

// Where where
func (g *gDb) Where(query interface{}, args ...interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Where(query, args...), g.transmit()...)
}

// Or or
func (g *gDb) Or(query interface{}, args ...interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Or(query, args...), g.transmit()...)
}

// Offset offset
func (g *gDb) Offset(offset int) standard.IBuilder {
	return CopyGDb(g.gDB().Offset(offset), g.transmit()...)
}

// Limit limit
func (g *gDb) Limit(limit int) standard.IBuilder {
	return CopyGDb(g.gDB().Limit(limit), g.transmit()...)
}

// Preload 预加载
func (g *gDb) Preload(query string, args ...interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Preload(query, args), g.transmit()...)
}

func (g *gDb) Select(query interface{}, args ...interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Select(query, args), g.transmit()...)
}

// Order order
func (g *gDb) Order(value interface{}) standard.IBuilder {
	var db = g.gDB()
	if reflect.ValueOf(value).Kind() == reflect.Slice {
		if sort, ok := value.([]interface{}); ok {
			for _, s := range sort {
				db = db.Order(s)
			}
		}
	} else {
		db = db.Order(value)
	}
	return CopyGDb(db, g.transmit()...)
}

// Begin begin
func (g *gDb) Begin(ctx context.Context) (standard.IBuilder, error) {
	var err error
	if g.TxCount() == 0 {
		g.db = g.gDB().Begin()
		err = g.gDB().Error
	}
	if nil == err {
		g.TxGagarin()
	}
	g.gDB().WithContext(ctx)
	return g, err
}

// Rollback rollback
func (g *gDb) Rollback() error {
	g.txClear()
	return g.gDB().Rollback().Error
}

// Commit commit
func (g *gDb) Commit() error {
	var err error
	if g.txRelief() == 0 {
		err = g.gDB().Commit().Error
	}
	return err
}

// Model model
func (g *gDb) Model(value interface{}) standard.IBuilder {
	return CopyGDb(g.gDB().Model(value), g.transmit()...)
}

// IsEmpty isEmpty
func (g *gDb) IsEmpty(err error) bool {
	return g.TransformNoRecordErr(err) == standard.ErrRecordNotFound
}

// IsWrite isWrite
func (g *gDb) IsWrite() bool {
	return g.Options.GetIsWrite()
}

// IsStartTx isStartTx
func (g *gDb) IsStartTx() bool {
	return g.Options.GetTxCount() > 0
}

// txRelief
func (g *gDb) txRelief() uint32 {
	g.Options.SetTxCount(g.Options.GetTxCount() - 1)
	return g.Options.GetTxCount()
}

// txClear
func (g *gDb) txClear() {
	g.Options.SetTxCount(0)
}

// TxGagarin txGagarin
func (g *gDb) TxGagarin() {
	g.Options.SetTxCount(g.Options.GetTxCount() + 1)
}

// TxCount txCount
func (g *gDb) TxCount() uint32 {
	return g.Options.GetTxCount()
}

// Exec exec
func (g *gDb) Exec(ctx context.Context, sql string, values ...interface{}) error {
	return g.gDB().WithContext(ctx).Exec(sql, values).Error
}

// Query query
func (g *gDb) Query(ctx context.Context, dest interface{}, sql string, values ...interface{}) error {
	return g.gDB().WithContext(ctx).Raw(sql, values...).Scan(dest).Error
}

// Copy copy
func (g *gDb) Copy() standard.IBuilder {
	return CopyGDb(g.gDB(), standard.WithIsWrite(g.IsWrite()), standard.WithTxCount(0))
}

// transmit
func (g *gDb) transmit() []standard.ServerOption {
	return []standard.ServerOption{
		standard.WithIsWrite(g.IsWrite()), standard.WithTxCount(g.TxCount()),
	}
}

func (g *gDb) TransformNoRecordErr(err error) error {
	if err != nil {
		if err == gorm.ErrRecordNotFound || err.Error() == "sql: no rows in result set" {
			err = standard.ErrRecordNotFound
		}
	}
	return err
}

func (g *gDb) GDB() *sql.DB {
	db, _ := g.gDB().DB()
	return db
}
