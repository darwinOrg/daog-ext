package daogext

import (
	"log"

	dgsys "github.com/darwinOrg/go-common/sys"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rolandhe/daog"
)

var (
	dataSource  daog.Datasource
	serviceName string
)

func SetDatasource(ds daog.Datasource) {
	dataSource = ds
}

func GetDatasource() daog.Datasource {
	return dataSource
}

// DbCfg 数据源配置, 包括数据库url和连接池相关配置
type DbCfg struct {
	// 数据库url
	Url string `json:"url" mapstructure:"url"`
	// 最大连接数
	MaxOpenConns int `json:"max-open-conns" mapstructure:"max-open-conns"`
	// 最大空闲连接数
	MaxIdleConns int `json:"max-idle-conns" mapstructure:"max-idle-conns"`
	// 连接的最大生命周期，单位是秒
	MaxLifeTime int `json:"max-life-time" mapstructure:"max-life-time"`
	// 最大空闲时间，单位是秒
	MaxIdleTime int `json:"max-idle-time" mapstructure:"max-idle-time"`
	// 不打印SQL日志
	NotLogSQL bool `json:"not-log-sql" mapstructure:"not-log-sql"`
}

func InitDbWithPossessionCallback(appName string, cfg *DbCfg) {
	InitDb(appName, cfg)

	daog.ChangeFieldOfInsBeforeWrite = func(valueMap map[string]any, extractor daog.FieldPointExtractor) error {
		return daog.ChangeInt64ByFieldNameCallback(valueMap, "op_id", extractor)
	}
	daog.AddNewModifyFieldBeforeUpdate = func(valueMap map[string]any, modifier daog.Modifier, existField func(filedName string) bool) error {
		return daog.ChangeModifierByFieldNameCallback(valueMap, "op_id", modifier, existField)
	}
}

func InitDb(appName string, cfg *DbCfg) {
	serviceName = appName
	dbConf := &daog.DbConf{
		DbUrl:    cfg.Url,
		Size:     cfg.MaxOpenConns,
		Life:     cfg.MaxLifeTime,
		IdleCons: cfg.MaxIdleConns,
		IdleTime: cfg.MaxIdleTime,
		LogSQL:   !cfg.NotLogSQL,
	}
	var err error
	dataSource, err = daog.NewDatasource(dbConf)
	if err != nil {
		if dgsys.IsFormalProfile() {
			panic(err)
		} else {
			log.Printf("init db error: %v", err)
		}
	}
}
