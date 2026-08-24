package daogext

import (
	"context"

	"github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

var CostThresholdMilli int64 = 500

func init() {
	daog.GLogger = &daogLogger{}
	dglogger.AppendIgnoreCallerFlags("daog(-ext)?(@[\\w.]+)?/([\\w.]+)?.go$")
}

type daogLogger struct {
}

func (dl *daogLogger) Error(ctx context.Context, err error) {
	alarmDatabaseError(getDgContext(ctx), err)
}

func (dl *daogLogger) Info(ctx context.Context, content string) {
}

func (dl *daogLogger) ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, sqlMd5 string) {
	dc := getDgContext(ctx)
	if sqlId := syncSqlContent(dc, sqlMd5, sql); sqlId > 0 {
		dglogger.Infof(dc, "%d | %s", sqlId, argsJson)
	} else {
		dglogger.Infof(dc, "%s | %s", sqlMd5, argsJson)
	}
}

func (dl *daogLogger) ExecSQLAfter(ctx context.Context, sqlMd5 string, cost int64) {
	if cost > CostThresholdMilli {
		dglogger.Infof(getDgContext(ctx), "%s | %dms", sqlMd5, cost)
	}
}

func (dl *daogLogger) SimpleLogError(err error) {
	alarmDatabaseError(dgctx.SimpleDgContext(), err)
}

var OnlyErrorLogger = &onlyErrorDaogLogger{}

type onlyErrorDaogLogger struct {
}

func (dl *onlyErrorDaogLogger) Error(ctx context.Context, err error) {
	alarmDatabaseError(getDgContext(ctx), err)
}

func (dl *onlyErrorDaogLogger) Info(ctx context.Context, content string) {
}

func (dl *onlyErrorDaogLogger) ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, sqlMd5 string) {
}

func (dl *onlyErrorDaogLogger) ExecSQLAfter(ctx context.Context, sqlMd5 string, cost int64) {
}

func (dl *onlyErrorDaogLogger) SimpleLogError(err error) {
	alarmDatabaseError(dgctx.SimpleDgContext(), err)
}

func getDgContext(ctx context.Context) *dgctx.DgContext {
	return &dgctx.DgContext{TraceId: daog.GetTraceIdFromContext(ctx), GoId: daog.GetGoroutineIdFromContext(ctx)}
}

func alarmDatabaseError(ctx *dgctx.DgContext, err error) {
	if errorProcessor != nil {
		errorProcessor(ctx, err)
	} else {
		dglogger.Errorf(ctx, "[daog] err: %v", err)
	}
}
