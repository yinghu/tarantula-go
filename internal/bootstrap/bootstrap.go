package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/metrics"
)

const (
	PUBLIC_ACCESS_CONTROL    int32 = 0
	PROTECTED_ACCESS_CONTROL int32 = 1
	ADMIN_ACCESS_CONTROL     int32 = 30
	SUDO_ACCESS_CONTROL      int32 = 100
)

const (
	DB_OP_ERR_CODE int = 500100

	WRONG_PASS_CODE int    = 400100
	WRONG_PASS_MSG  string = "wrong user/password"

	ILLEGAL_TOKEN_CODE  int = 400101
	INVALID_TOKEN_CODE  int = 400102
	ILLEGAL_ACCESS_CODE int = 400103

	INVALID_TOKEN_MSG string = "invalid token"
)

const (
	METRICS_SQL_SCHEMA         string = "CREATE TABLE IF NOT EXISTS req_metrics (id BIGSERIAL PRIMARY KEY,path VARCHAR(50) NOT NULL,req_timed BIGINT NOT NULL,req_time TIMESTAMP DEFAULT NOW(),node VARCHAR(10) NOT NULL,req_id INTEGER DEFAULT 0,req_code INTEGER DEFAULT 0)"
	ITEM_ENUM_SQL_SCHEMA       string = "CREATE TABLE IF NOT EXISTS item_enum (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE)"
	ITEM_ENUM_VALUE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_enum_value (enum_id INTEGER,name VARCHAR(100) NOT NULL,value INTEGER NOT NULL,PRIMARY KEY(enum_id,name))"
	ITEM_CATEGORY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE,scope VARCHAR(30) NOT NULL ,rechargeable BOOL NOT NULL ,description VARCHAR(100) NOT NULL)"
	ITEM_PROPERTY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_property (category_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(100) NOT NULL ,reference VARCHAR(100) NOT NULL ,nullable BOOL NOT NULL ,downloadable BOOL NOT NULL, PRIMARY KEY(category_id,name))"

	ITEM_CONFIGURATION_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_configuration (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL,type VARCHAR(50) NOT NULL ,type_id VARCHAR(50) NOT NULL ,category VARCHAR(100) NOT NULL ,version VARCHAR(10) NOT NULL,UNIQUE(name,version))"
	ITEM_HEADER_SQL_SCHEMA        string = "CREATE TABLE IF NOT EXISTS item_header (configuration_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,value VARCHAR(100) NOT NULL, PRIMARY KEY(configuration_id,name))"
	ITEM_APPLICATION_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_application (configuration_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,reference_id BIGINT NOT NULL,PRIMARY KEY(configuration_id,name,reference_id))"
)

type TarantulaContext interface {
	Config() string
	Start(f conf.Env, c cluster.Cluster) error
	Shutdown()
	event.EventService
	cluster.KeyListener
}

type TarantulaApp interface {
	ItemService() item.ItemService
	Metrics() metrics.MetricsService
	Cluster() cluster.Cluster
	Authenticator() core.Authenticator
	AccessControl() int32
	Request(sesion core.OnSession, w http.ResponseWriter, r *http.Request)
}
