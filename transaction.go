package daogext

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
)

var dataSource daog.Datasource

func InitDatasource(ds daog.Datasource) {
	dataSource = ds
}

func ReadonlyWithResult[T any](ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) (T, error)) (T, error) {
	return daog.AutoTransWithResult(func() (*daog.TransContext, error) {
		return daog.NewTransContext(dataSource, txrequest.RequestReadonly, ctx.TraceId)
	}, workFn)
}

func Write(ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) error) error {
	return daog.AutoTrans(func() (*daog.TransContext, error) {
		return daog.NewTransContext(dataSource, txrequest.RequestWrite, ctx.TraceId)
	}, workFn)
}

func WriteWithResult[T any](ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) (T, error)) (T, error) {
	return daog.AutoTransWithResult(func() (*daog.TransContext, error) {
		return daog.NewTransContext(dataSource, txrequest.RequestWrite, ctx.TraceId)
	}, workFn)
}
