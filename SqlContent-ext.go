package daogext

import (
	"sync"
	"time"

	"github.com/darwinOrg/go-common/context"
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var SqlContentExtDao = &sqlContentExtDao{}

type sqlContentExtDao struct{}

func (d *sqlContentExtDao) FindByCategory(ctx *dgctx.DgContext, tc *daog.TransContext, category string) ([]*SqlContent, error) {
	list, err := SqlContentDao.QueryListMatcherWithViewColumns(tc, daog.NewMatcher().Eq(SqlContentFields.Category, category),
		[]string{SqlContentFields.Id, SqlContentFields.SqlMd5})
	if err != nil {
		dglogger.Errorf(ctx, "SqlContentDao.GetById error: %v", err)
		return nil, err
	}

	return list, nil
}

func (d *sqlContentExtDao) GetByCategoryAndSqlMd5(ctx *dgctx.DgContext, tc *daog.TransContext, category, sqlMd5 string) (*SqlContent, error) {
	sc, err := SqlContentDao.QueryOneMatcher(tc, daog.NewMatcher().Eq(SqlContentFields.Category, category).Eq(SqlContentFields.SqlMd5, sqlMd5))
	if err != nil {
		dglogger.Errorf(ctx, "SqlContentDao.QueryOneMatcher error: %v", err)
		return nil, err
	}

	return sc, nil
}

func (d *sqlContentExtDao) Create(ctx *dgctx.DgContext, tc *daog.TransContext, sc *SqlContent) error {
	now := ttypes.NormalDatetime(time.Now())
	sc.CreatedAt = now

	_, err := SqlContentDao.Insert(tc, sc)
	if err != nil {
		dglogger.Errorf(ctx, "SqlContentDao.Insert error: %v", err)
		return err
	}

	return nil
}

var sqlContentMap = sync.Map{}

func initSqlContent() {
	ctx := dgctx.SimpleDgContext()

	scList, err := ReadonlyWithResult(ctx, func(tc *daog.TransContext) ([]*SqlContent, error) {
		return SqlContentExtDao.FindByCategory(ctx, tc, dgsys.ServiceName)
	})
	if err != nil {
		dglogger.Warnf(ctx, "SqlContentExtDao.FindByCategory err: %v", err)
		return
	}
	if len(scList) == 0 {
		return
	}

	for _, sc := range scList {
		sqlContentMap.Store(sc.SqlMd5, sc.Id)
	}
}

func syncSqlContent(ctx *dgctx.DgContext, sqlMd5, content string) (int64, error) {
	if sqlId, ok := sqlContentMap.Load(sqlMd5); ok {
		return sqlId.(int64), nil
	}

	return WriteWithResult(ctx, func(tc *daog.TransContext) (int64, error) {
		sc, err := SqlContentExtDao.GetByCategoryAndSqlMd5(ctx, tc, dgsys.ServiceName, sqlMd5)
		if err != nil {
			return 0, err
		}
		if sc != nil {
			sqlContentMap.Store(sc.SqlMd5, sc.Id)
			return sc.Id, nil
		}

		sc = &SqlContent{
			Category: dgsys.ServiceName,
			SqlMd5:   sqlMd5,
			Content:  content,
		}
		err = SqlContentExtDao.Create(ctx, tc, sc)
		if err != nil {
			return 0, err
		}

		sqlContentMap.Store(sc.SqlMd5, sc.Id)
		return sc.Id, nil
	})
}
