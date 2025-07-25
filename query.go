package daogext

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

func CountRaw(ctx *dgctx.DgContext, sql string, args ...any) (int64, error) {
	return ReadonlyWithResult(ctx, func(tc *daog.TransContext) (int64, error) {
		return CountRawByTc(ctx, tc, sql, args...)
	})
}

func CountRawByTc(ctx *dgctx.DgContext, tc *daog.TransContext, sql string, args ...any) (int64, error) {
	scs, err := daog.QueryRawSQL(tc, func(ins *int64) []any {
		return []any{ins}
	}, sql, args...)
	if err != nil {
		dglogger.Errorf(ctx, "daog.QueryRawSQL error: %v", err)
		return 0, dgerr.SYSTEM_ERROR
	}
	if len(scs) == 0 {
		return 0, nil
	}

	return *scs[0], nil
}

func QueryRawList[T any](ctx *dgctx.DgContext, sql string, args ...any) ([]*T, error) {
	return ReadonlyWithResult(ctx, func(tc *daog.TransContext) ([]*T, error) {
		return QueryRawListByTc[T](ctx, tc, sql, args...)
	})
}

func QueryRawListByTc[T any](ctx *dgctx.DgContext, tc *daog.TransContext, sql string, args ...any) ([]*T, error) {
	list, err := daog.QueryRawSQL(tc, func(ins *T) []any {
		return utils.ReflectAllFieldValuePointers(ins)
	}, sql, args...)
	if err != nil {
		dglogger.Errorf(ctx, "daog.QueryRawSQL error: %v", err)
	}
	return list, err
}
