package dbx

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/thoohv5/orm/standard"
)

var (
	dbWrapperList sync.Map
)

type Dao struct {
	*CommonReq
	tabler  Tabler
	filter  IFilter
	builder standard.IBuilder
}

// NewDao new
func NewDao(model IDao) *Dao {
	// connect, ok := model.(Connecter)
	// if !ok {
	// 	return nil
	// }

	table, ok := model.(Tabler)
	if !ok {
		return nil
	}
	filter, ok := model.(IFilter)
	if !ok {
		return nil
	}
	return &Dao{
		tabler: table,
		filter: filter,
		// builder: connect.Connection(),
	}
}

// RegisterDao register dao
func RegisterDao(model IDao) *Dao {
	return NewDao(model)
}

func (d *Dao) Add(ctx context.Context, value Tabler) error {
	return d.builder.Write().Add(ctx, value)
}

func (d *Dao) Update(ctx context.Context, param IQuery, update map[string]interface{}) error {
	if err := d.check(param); nil != err {
		return err
	}
	return d.filter.BuildFilterQuery(d.builder.Write(), param).Model(param).Update(ctx, update)
}

func (d *Dao) Delete(ctx context.Context, param IQuery) error {
	if err := d.check(param); nil != err {
		return err
	}
	return d.filter.BuildFilterQuery(d.builder.Write(), param).Delete(ctx, d.tabler)
}

func (d *Dao) Find(ctx context.Context, param IQuery, result Tabler) error {
	if err := d.check(param); nil != err {
		return err
	}
	return d.filter.Filter(d.filter.BuildFilterQuery(d.builder.Read(), param), param.GetCommonReq()).Find(ctx, result)
}

func (d *Dao) Get(ctx context.Context, param IQuery, result interface{}) error {
	if err := d.check(param); nil != err {
		return err
	}
	return d.filter.Filter(d.filter.BuildFilterQuery(d.builder.Read(), param), param.GetCommonReq()).Get(ctx, result)
}

func (d *Dao) GetCommonReq() *CommonReq {
	if nil == d {
		return &CommonReq{}
	}
	return d.CommonReq
}

// filter
func (d *Dao) Filter(build standard.IBuilder, condition *CommonReq) standard.IBuilder {
	if nil != condition {
		// Start
		if start := condition.Start; start > 0 {
			build = build.Offset(start)
		}
		// Limit
		limit := condition.Limit
		if limit == 0 {
			limit = 500
		}
		build = build.Limit(limit)
		// 排序
		for _, sort := range condition.Sorts {
			if string(sort[0]) == "+" {
				sort = sort[1:] + " ASC"
			} else if string(sort[0]) == "-" {
				sort = sort[1:] + " DESC"
			} else {
				sort = sort + " ASC"
			}
			build = build.Order(sort)
		}
	}

	return build
}

// builder
func (d *Dao) GetBuilder() standard.IBuilder {
	return d.builder
}

// builder
func (d *Dao) SetBuilder(builder standard.IBuilder) {
	d.builder = builder
}

func (d *Dao) check(dest interface{}) error {
	// 校验类型
	dv := reflect.ValueOf(dest)
	dt := dv.Type()
	if dt.Kind() != reflect.Ptr {
		return errors.New("被赋值的单体必须是指针类型")
	}
	return nil
}
