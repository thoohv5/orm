package dbx

import (
	"context"
	"fmt"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"

	"github.com/thoohv5/orm/standard"
)

const (
	intErrCode = 1
	dbErrCode  = 2

	transErrF = "TransError[%d]: %s"
)

type TransError struct {
	Code int
	Msg  string
}

func (t TransError) Error() string {
	return fmt.Sprintf(transErrF, t.Code, t.Msg)
}

type Task func(builder standard.IBuilder) error

func errHandler(builder standard.IBuilder, task Task) (err error) {
	defer func() {
		if e := recover(); e != nil {
			msg := fmt.Sprintf("panic: %s\ncalltrace : %s", fmt.Sprint(e), string(debug.Stack()))
			err = &TransError{intErrCode, msg}
		}
	}()
	return task(builder)
}

func ExecTrans(ctx context.Context, builder standard.IBuilder, trans ...Task) error {
	dber, err := builder.Begin(ctx)
	if err != nil {
		fmt.Printf("gDB begin transaction failed: %s", err.Error())
		return &TransError{dbErrCode, err.Error()}
	}
	for _, task := range trans {
		if err := errHandler(dber, task); err != nil {
			if err := dber.Rollback(); err != nil {
				fmt.Printf("roll_back : %s", err.Error())
			}
			return err
		}
	}

	if err := dber.Commit(); err != nil {
		_ = dber.Rollback()
		return &TransError{dbErrCode, fmt.Sprintf("%v", err.Error())}
	}
	return nil
}
