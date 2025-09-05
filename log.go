package daogext

import (
	"context"
	"fmt"

	alarmsdk "e.globalpand.cn/libs/alarm-sdk"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

func init() {
	daog.GLogger = &daogLogger{}
}

type daogLogger struct {
}

func (dl *daogLogger) Error(ctx context.Context, err error) {
	alarmDatabaseError(dgctx.SimpleDgContext(), err)
}

func (dl *daogLogger) Info(ctx context.Context, content string) {
	dglogger.Infof(getDgContext(ctx), "[daog] content: %s", content)
}

func (dl *daogLogger) ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, sqlMd5 string) {
	dglogger.Infof(getDgContext(ctx), "[daog] [Trace SQL] sqlMd5=%s, sql: %s, args:%s", sqlMd5, sql, argsJson)
}

func (dl *daogLogger) ExecSQLAfter(ctx context.Context, sqlMd5 string, cost int64) {
	dglogger.Infof(getDgContext(ctx), "[daog] [Trace SQL] sqlMd5=%s, cost %d ms", sqlMd5, cost)
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
	alarmsdk.BackendAlarm(ctx, fmt.Sprintf("database execution error: %v", err))
}
