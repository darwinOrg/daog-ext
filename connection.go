package daogext

import "github.com/rolandhe/daog"

var dataSource daog.Datasource

func SetDatasource(ds daog.Datasource) {
	dataSource = ds
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
}

func InitDb(dc *DbCfg) {
	dbConf := &daog.DbConf{
		DbUrl:    dc.Url,
		Size:     dc.MaxOpenConns,
		Life:     dc.MaxLifeTime,
		IdleCons: dc.MaxIdleConns,
		IdleTime: dc.MaxIdleTime,
		LogSQL:   true,
	}
	var err error
	dataSource, err = daog.NewDatasource(dbConf)
	if err != nil {
		panic(err)
	}
}
