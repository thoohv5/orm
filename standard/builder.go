package standard

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// IBuilder 标准
type IBuilder interface {
	// Add 添加
	Add(ctx context.Context, value interface{}) error
	// Update 更新
	Update(ctx context.Context, attrs interface{}) error
	// Delete 删除
	Delete(ctx context.Context, value interface{}, where ...interface{}) error
	// Find 查询单条
	Find(ctx context.Context, out interface{}, where ...interface{}) error
	// Get 查询多条
	Get(ctx context.Context, out interface{}, where ...interface{}) error
	GetGDB() *gorm.DB
	// Count 统计
	Count(ctx context.Context, value *int64) error
	// Where 查询且条件
	Where(query interface{}, args ...interface{}) IBuilder
	// Set 设置属性
	Set(name string, value interface{}) IBuilder
	// Or 查询或条件
	Or(query interface{}, args ...interface{}) IBuilder
	// Offset 查询定位
	Offset(offset int) IBuilder
	// Limit 查询区间
	Limit(limit int) IBuilder
	// Order 排序
	Order(value interface{}) IBuilder
	// Begin 事务开始
	Begin(ctx context.Context) (IBuilder, error)
	// Preload 预加载
	Preload(query string, args ...interface{}) IBuilder
	Select(query interface{}, args ...interface{}) IBuilder
	// Session
	Session(sess *gorm.Session) IBuilder
	// Rollback 事务回滚
	Rollback() error
	// Commit 事务提交
	Commit() error
	// Model 模型
	Model(value interface{}) IBuilder
	// IsEmpty 检查数据是否为空
	IsEmpty(e error) bool
	// IsWrite 检查是否为写库
	IsWrite() bool
	// IsStartTx 检查是否在事物中
	IsStartTx() bool
	// Exec 执行
	Exec(ctx context.Context, sql string, values ...interface{}) error
	// Query 查询
	Query(ctx context.Context, dest interface{}, sql string, values ...interface{}) error
	// Copy 复制
	Copy() IBuilder

	Write() IBuilder

	Read() IBuilder

	GDB() *sql.DB
}
