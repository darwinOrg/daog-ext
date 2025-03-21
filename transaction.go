package daogext

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
)

func Readonly(ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) error) error {
	return daog.AutoTrans(func() (*daog.TransContext, error) {
		return newReadonlyTransContext(ctx)
	}, workFn)
}

func ReadonlyWithResult[T any](ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) (T, error)) (T, error) {
	return daog.AutoTransWithResult(func() (*daog.TransContext, error) {
		return newReadonlyTransContext(ctx)
	}, workFn)
}

func Write(ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) error) error {
	return daog.AutoTrans(func() (*daog.TransContext, error) {
		return newWriteTransContext(ctx)
	}, workFn)
}

func WriteWithResult[T any](ctx *dgctx.DgContext, workFn func(tc *daog.TransContext) (T, error)) (T, error) {
	return daog.AutoTransWithResult(func() (*daog.TransContext, error) {
		return newWriteTransContext(ctx)
	}, workFn)
}

func newReadonlyTransContext(ctx *dgctx.DgContext) (*daog.TransContext, error) {
	return newTransContextWithRequestStyle(ctx, txrequest.RequestReadonly)
}

func newWriteTransContext(ctx *dgctx.DgContext) (*daog.TransContext, error) {
	tc, err := newTransContextWithRequestStyle(ctx, txrequest.RequestWrite)
	if tc != nil {
		tc.ExtInfo = map[string]any{
			"op_id": ctx.OpId,
		}
	}
	return tc, err
}

func newTransContextWithRequestStyle(ctx *dgctx.DgContext, requestStyle txrequest.RequestStyle) (*daog.TransContext, error) {
	tc, err := daog.NewTransContext(dataSource, requestStyle, ctx.TraceId)
	if err != nil {
		return tc, err
	}

	if tc.LogSQL && ctx.NotLogSQL {
		tc.LogSQL = false
	}

	return tc, err
}
