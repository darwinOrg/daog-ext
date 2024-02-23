package daogext

import (
	"context"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/google/uuid"
	"github.com/rolandhe/daog"
)

var OnlyErrorLogger = &onlyErrorDaogLogger{}

func init() {
	daog.GLogger = &daogLogger{}
}

type daogLogger struct {
}

func (dl *daogLogger) Error(ctx context.Context, err error) {
	dglogger.Errorf(getDgContext(ctx), "[daog] err: %v", err)
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
	dglogger.Errorf(&dgctx.DgContext{TraceId: uuid.NewString()}, "[daog] err: %v", err)
}

type onlyErrorDaogLogger struct {
}

func (dl *onlyErrorDaogLogger) Error(ctx context.Context, err error) {
	dglogger.Errorf(getDgContext(ctx), "[daog] err: %v", err)
}

func (dl *onlyErrorDaogLogger) Info(ctx context.Context, content string) {
}

func (dl *onlyErrorDaogLogger) ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, sqlMd5 string) {
}

func (dl *onlyErrorDaogLogger) ExecSQLAfter(ctx context.Context, sqlMd5 string, cost int64) {
}

func (dl *onlyErrorDaogLogger) SimpleLogError(err error) {
	dglogger.Errorf(&dgctx.DgContext{TraceId: uuid.NewString()}, "[daog] err: %v", err)
}

func getDgContext(ctx context.Context) *dgctx.DgContext {
	return &dgctx.DgContext{TraceId: daog.GetTraceIdFromContext(ctx), GoId: daog.GetGoroutineIdFromContext(ctx)}
}
