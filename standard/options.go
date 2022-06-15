package standard

import (
	"sync"

	"github.com/thoohv5/logger"
)

// Options 可选参数列表
type Options struct {
	isWrite bool
	txCount uint32

	isSetLog bool
	log      logger.ILogger
	level    string

	syncRWMutex sync.RWMutex
}

// ServerOption 为可选参数赋值的函数
type ServerOption func(*Options)

// WithIsWrite 是否为写库
func WithIsWrite(isWrite bool) ServerOption {
	return func(o *Options) {
		o.isWrite = isWrite
	}
}

func WithLogger(log logger.ILogger) ServerOption {
	return func(o *Options) {
		o.isSetLog = true
		o.log = log
	}
}

func WithLoggerLevel(level string) ServerOption {
	return func(o *Options) {
		o.level = level
	}
}

func (o *Options) GetIsSetLog() bool {
	return o.isSetLog
}

func (o *Options) GetLogger() logger.ILogger {
	return o.log
}

func (o *Options) GetLoggerLevel() string {
	return o.level
}

// GetIsWrite 获取是否为写库
func (o *Options) GetIsWrite() bool {
	return o.isWrite
}

// WithTxCount 事物统计
func WithTxCount(count uint32) ServerOption {
	return func(o *Options) {
		o.txCount = count
	}
}

// GetTxCount 获取事物统计
func (o *Options) GetTxCount() uint32 {
	defer o.syncRWMutex.RUnlock()
	o.syncRWMutex.RLock()
	return o.txCount
}

// SetTxCount 设置事物统计
func (o *Options) SetTxCount(txCount uint32) {
	defer o.syncRWMutex.Unlock()
	o.syncRWMutex.Lock()
	o.txCount = txCount
}
