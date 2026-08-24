package daogext

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var SqlContentFields = struct {
	Id        string
	Category  string
	SqlMd5    string
	Content   string
	CreatedAt string
}{
	"id",
	"category",
	"sql_md5",
	"content",
	"created_at",
}

var SqlContentMeta = &daog.TableMeta[SqlContent]{
	Table: "sql_content",
	Columns: []string{
		"id",
		"category",
		"sql_md5",
		"content",
		"created_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *SqlContent, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "category" == columnName {
			if point {
				return &ins.Category
			}
			return ins.Category
		}
		if "sql_md5" == columnName {
			if point {
				return &ins.SqlMd5
			}
			return ins.SqlMd5
		}
		if "content" == columnName {
			if point {
				return &ins.Content
			}
			return ins.Content
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}

		return nil
	},
}

var SqlContentDao daog.QuickDao[SqlContent] = &struct {
	daog.QuickDao[SqlContent]
}{
	NewBaseQuickDao(SqlContentMeta),
}

type SqlContent struct {
	Id        int64                 `json:"id" remark:"id"`
	Category  string                `json:"category" remark:"类别"`
	SqlMd5    string                `json:"sqlMd5" remark:"sql md5"`
	Content   string                `json:"content" remark:"内容"`
	CreatedAt ttypes.NormalDatetime `json:"createdAt" remark:"创建时间"`
}
