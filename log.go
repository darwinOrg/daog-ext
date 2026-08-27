package daogext

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/utils"
	"github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

var CostThresholdMilli int64 = 500

// 多 ? 占位符严格匹配：括号内是 ? 用逗号分隔，各部分之间允许任意空格
var multiPlaceholderRe = regexp.MustCompile(`\(\s*\?\s*(?:,\s*\?\s*)*\)`)

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

func (dl *daogLogger) ExecSQLBefore(ctx context.Context, sql string, argsJson []byte, _ string) {
	sql = strings.TrimSpace(sql)
	sql = multiPlaceholderRe.ReplaceAllString(sql, "(?)")
	sqlMd5 := utils.Md5Hex(sql)
	if len(argsJson) > 0 {
		argsJson = bytes.TrimSuffix(bytes.TrimPrefix(argsJson, []byte("[")), []byte("]"))
	}
	dc := getDgContext(ctx)
	if sqlId := syncSqlContent(dc, sqlMd5, sql); sqlId > 0 {
		if len(argsJson) > 0 {
			dglogger.Infof(dc, "%d | %s", sqlId, argsJson)
		} else {
			dglogger.Infof(dc, "%d", sqlId)
		}
	} else if len(argsJson) > 0 {
		dglogger.Infof(dc, "%s | %s", sqlMd5, argsJson)
	} else {
		dglogger.Infof(dc, "%s", sqlMd5)
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
