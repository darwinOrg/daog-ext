package daogext

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/darwinOrg/go-common/collection"
	"github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/sys"
	"github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
)

type columnInfo struct {
	TableName  string
	ColumnName string
	ColumnType string
}

type tableMetaExt struct {
	Table   string
	Columns []string
	Types   []string
}

var tableMetaExts []*tableMetaExt

func NewBaseQuickDao[T any](meta *daog.TableMeta[T]) daog.QuickDao[T] {
	tableMetaExts = append(tableMetaExts, convertToMetaExt(meta))
	return daog.NewBaseQuickDao(meta)
}

func convertToMetaExt[T any](meta *daog.TableMeta[T]) *tableMetaExt {
	obj := new(T)
	var types []string
	for _, column := range meta.Columns {
		fieldObj := meta.LookupFieldFunc(column, obj, true)
		fieldType := reflect.TypeOf(fieldObj)
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		fieldTypeString := fieldType.String()
		types = append(types, fieldTypeString)
	}

	return &tableMetaExt{
		Table:   meta.Table,
		Columns: meta.Columns,
		Types:   types,
	}
}

func validateTableMeta() {
	if !dgsys.IsFormalProfile() || len(tableMetaExts) == 0 {
		return
	}

	tableMetaExtMap := dgcoll.Trans2Map(tableMetaExts, func(meta *tableMetaExt) string { return meta.Table })
	tablesNames := dgcoll.MapToList(tableMetaExts, func(meta *tableMetaExt) string { return meta.Table })
	tablesNamesStr := "'" + strings.Join(tablesNames, "','") + "'"
	queryColumnSql := fmt.Sprintf(`
		SELECT  TABLE_NAME, COLUMN_NAME, COLUMN_TYPE
        FROM INFORMATION_SCHEMA.COLUMNS 
        WHERE TABLE_SCHEMA = DATABASE() 
          AND TABLE_NAME in (%s)
        ORDER BY TABLE_NAME, ORDINAL_POSITION
	`, tablesNamesStr)

	ctx := dgctx.SimpleDgContext()

	columnInfos, err := ReadonlyWithResult(ctx, func(tc *daog.TransContext) ([]*columnInfo, error) {
		list, err := daog.QueryRawSQL(tc, func(ins *columnInfo) []any {
			return []any{&ins.TableName, &ins.ColumnName, &ins.ColumnType}
		}, queryColumnSql)
		if err != nil {
			dglogger.Errorf(ctx, "查询数据库列元信息错误: %v", err)
			return nil, err
		}
		return list, nil
	})
	if err != nil {
		return
	}

	existsTableNames := dgcoll.MapToSet(columnInfos, func(info *columnInfo) string { return info.TableName })
	notExistsTableNames := dgcoll.Remove(tablesNames, existsTableNames)
	if len(notExistsTableNames) > 0 {
		dbe := fmt.Errorf("错误！数据库缺少表: %s", strings.Join(notExistsTableNames, ", "))
		if errorProcessor != nil {
			errorProcessor(ctx, dbe)
		} else {
			dglogger.Warn(ctx, dbe)
		}
	}

	tableName2ColumnInfosMap := dgcoll.GroupBy(columnInfos, func(info *columnInfo) string { return info.TableName })
	for tableName, tableColumnInfos := range tableName2ColumnInfosMap {
		metaExt := tableMetaExtMap[tableName]
		metaColumns := metaExt.Columns
		var missColumns []string
		var disMatchColumns []string

		for i, metaColumn := range metaColumns {
			tableColumnInfo := dgcoll.FindFirst(tableColumnInfos, func(info *columnInfo) bool {
				return strings.ReplaceAll(info.ColumnName, "`", "") == strings.ReplaceAll(metaColumn, "`", "")
			}, nil)

			// 如果实际数据库里面没有这个字段，则报警
			if tableColumnInfo == nil {
				missColumns = append(missColumns, metaColumn)
				continue
			}

			metaColumnType := metaExt.Types[i]
			dbColumnType := tableColumnInfo.ColumnType

			// 如果mysql与go的数据类型不匹配，则报警
			if !isMySQLTypeCompatibleWithGo(dbColumnType, metaColumnType) {
				disMatchColumns = append(disMatchColumns, metaColumn)
				continue
			}
		}

		if len(missColumns) > 0 {
			dbe := fmt.Errorf("错误！[%s - %s]字段缺失", tableName, strings.Join(missColumns, ", "))
			if errorProcessor != nil {
				errorProcessor(ctx, dbe)
			} else {
				dglogger.Warn(ctx, dbe)
			}
		}

		if len(disMatchColumns) > 0 {
			dbe := fmt.Errorf("错误！[%s - %s]字段类型不匹配", tableName, strings.Join(disMatchColumns, ", "))
			if errorProcessor != nil {
				errorProcessor(ctx, dbe)
			} else {
				dglogger.Warn(ctx, dbe)
			}
		}
	}
}
