package qsmysql

import (
	"math/rand"
	//Register mysql import
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	defaultMaxConns     = 256
	defaultMaxIdleConns = 128
	defaultMaxLifetime  = 5 // in minutes
	defaultMySQLPOrt    = 3306
	defaultCharset      = "utf8mb4"
	defaultLogMode      = false
	defaultDriverName   = "mysql"
	once                sync.Once
	qsmysql             *QSMySQL
)

type mysqlMaster struct {
	handler *gorm.DB
	dsn     string
	conf    connPoolConf
}

type mysqlSlave struct {
	handlers []*gorm.DB
	dsns     []string
	conf     connPoolConf
}

type connPoolConf struct {
	maxOpenConns, maxIdleConns, connMaxLifetime int
	logMode                                     bool
}

type QSMySQL struct {
	master    *mysqlMaster
	slave     *mysqlSlave
	viperConf *viper.Viper
	hasslave  bool
}

func init() {
	qsmysql = New()
	// qsmysql.master.handler = new(gorm.DB)
	// qsmysql.slave.handlers = make([]*gorm.DB, 0, 8)
}

func New() *QSMySQL {
	return &QSMySQL{}
}

func SetConfig(v *viper.Viper) error {
	return qsmysql.SetConfig(v)
}

func (q *QSMySQL) SetConfig(v *viper.Viper) error {
	q.viperConf = v
	if v.IsSet("master") {
		dsns, confOrm, err := parseViper(v.Sub("master"))
		if err != nil {
			//todo: more readable log
			return err
		}
		q.master = new(mysqlMaster)
		q.master.dsn = dsns[0]
		q.master.conf = confOrm
	}
	if v.IsSet("slave") {
		dsns, confOrm, err := parseViper(v.Sub("slave"))
		if err != nil {
			return err
		}
		q.slave = new(mysqlSlave)
		q.hasslave = true
		q.slave.dsns = dsns
		q.slave.conf = confOrm
	}
	return nil
}

func Close() error {
	return qsmysql.Close()
}

func (q *QSMySQL) Close() error {
	if q.master.handler != nil {
		//TODO: handle the error that Close() returned
		q.master.handler.Close()
	}
	if q.hasslave {
		for _, h := range q.slave.handlers {
			h.Close()
		}
	}
	q.viperConf = nil
	return nil
}

/**
mysql:
  master:
    host: rm-bp14040c8no99686c.mysql.rds.aliyuncs.com
    port: 3306
    user: youth
    password: borA6i@#$%^&teuwtMq6F9eYnY
    db: youth
    charset: utf8mb4
    max_idle_conns: 20
    max_open_conns: 200
    log_mode: false
    ## minutes
	conn_max_lifetime: 6
  slave:
	host:
		-rm-bp14040c8no99686c.mysql.rds.aliyuncs.com
		-rm-bp14040c8no99686c.mysql.rds.aliyun.com
    port: 3306
    user: youth
    password: borA6i@#$%^&teuwtMq6F9eYnY
    db: youth
    charset: utf8mb4
    max_idle_conns: 20
    max_open_conns: 200
    log_mode: false
    ## minutes
	conn_max_lifetime: 6

**/
func parseViper(v *viper.Viper) (dsns []string, conf connPoolConf, err error) {
	var (
		hc             = 1
		host           string
		hosts          []string
		port           int
		user, password string
		dbname         string
		charset        string
	)

	if v.IsSet("host") {
		h := v.Get("host")
		switch h.(type) {
		case string:
			host = h.(string)
		case interface{}:
			hosts = cast.ToStringSlice(h)
			hc = len(host)
		default:
			panic(fmt.Errorf("got host type error %s,%v",
				reflect.TypeOf(h), reflect.ValueOf(h)))
		}
	} else {
		panic(fmt.Errorf("invalid host specified"))
	}
	if !v.IsSet("port") {
		port = defaultMySQLPOrt
	} else {
		port = v.GetInt("port")
	}
	if !v.IsSet("user") || !v.IsSet("password") {
		panic(fmt.Errorf("invalid user specified"))
	} else {
		user = v.GetString("user")
		password = v.GetString("password")
	}
	if !v.IsSet("db") {
		panic("invalid user specified")
	} else {
		dbname = v.GetString("db")
	}
	if !v.IsSet("charset") {
		charset = defaultCharset
	}
	if !v.IsSet("max_open_conns") {
		conf.maxOpenConns = defaultMaxConns
	} else {
		conf.maxOpenConns = v.GetInt("max_open_conns")
	}
	if !v.IsSet("max_idle_conns") {
		conf.maxIdleConns = defaultMaxIdleConns
	} else {
		conf.maxIdleConns = v.GetInt("max_idle_conns")
	}
	if !v.IsSet("conn_max_lifetime") {
		conf.connMaxLifetime = defaultMaxLifetime
	} else {
		conf.connMaxLifetime = v.GetInt("conn_max_lifetime")
	}
	if !v.IsSet("log_mode") {
		conf.logMode = defaultLogMode
	} else {
		conf.logMode = v.GetBool("log_mode")
	}

	if !v.IsSet("charset") {
		charset = defaultCharset
	} else {
		charset = v.GetString("charset")
	}

	if hc == 1 {
		d := user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=" + charset + "&parseTime=True&loc=Local"
		dsns = append(dsns, d)
	} else {
		for _, h := range hosts {
			d := user + ":" + password + "@tcp(" + h + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=" + charset + "&parseTime=True&loc=Local"
			dsns = append(dsns, d)
		}
	}
	return dsns, conf, nil
}

func ConnDB(dsn string, conf connPoolConf) (*gorm.DB, error) {
	db, err := gorm.Open(defaultDriverName, dsn)
	if err != nil {
		//TODO: Log error
		return nil, err
	}
	db.DB().SetMaxOpenConns(conf.maxOpenConns)
	db.DB().SetMaxIdleConns(conf.maxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(conf.connMaxLifetime) * time.Minute)
	db.LogMode(conf.logMode)
	return db, nil
}

func (q *QSMySQL) connMaster() error {
	db, err := ConnDB(q.master.dsn, q.master.conf)
	if err != nil {
		return err
	}
	q.master.handler = db
	return nil
}

func (q *QSMySQL) connSlave() error {
	q.slave.handlers = make([]*gorm.DB, len(q.slave.dsns))
	for _, d := range q.slave.dsns {
		db, err := ConnDB(d, q.slave.conf)
		if err != nil {
			return err
		} else {
			q.slave.handlers = append(q.slave.handlers, db)
		}
	}
	return nil
}

func GetMaster() *gorm.DB {
	return qsmysql.GetMaster()
}

func (q *QSMySQL) GetMaster() *gorm.DB {
	live := true
	if q.master.handler == nil {
		live = false
	}
	// if err := qsmysql.master.handler.DB().Ping(); err != nil {
	// 	live = false
	// }

	if live == false {
		once.Do(func() {
			if err := q.connMaster(); err != nil {
				panic("connMaster failed, error: " + err.Error())
			}
		})
	}
	//TODO: attention!!! here may cause panic,while users use chained call like qsmysql.GetMaster().Query() etc.
	return q.master.handler
}

func GetSlave() *gorm.DB {
	return qsmysql.GetSlave()
}

func (q *QSMySQL) GetSlave() *gorm.DB {
	if !q.hasslave {
		return q.GetMaster()
	}

	live := true
	if q.slave == nil {
		live = false
	}
	if !live {
		once.Do(func() {
			if err := q.connSlave(); err != nil {
				panic("connSlave failed: " + err.Error())
			}
		})
	}
	rand.Seed(time.Now().UTC().UnixNano())
	//TODO: attention!! here may cause panic, while users use chained call like qsmysql.GetSlave.Query() etc.
	return q.slave.handlers[rand.Intn(len(qsmysql.slave.handlers)+1)]
}

//todo: exportStats should export the stats of *grom.DB within a seprated goroutine
func (q *QSMySQL) exportStats() {

	if q.master.handler != nil {

	}
	if q.hasslave && len(q.slave.handlers) > 0 {

	}

}

func getStats(g *gorm.DB) {
	_ = g.DB().Stats

}

func SetRemote(remoteProvider string, endpoint, path string) error {
	if remoteProvider != "etcd" {
		return errors.New("remote provider " + remoteProvider + " is not supported")
	}
	v := viper.New()
	err := v.AddRemoteProvider(remoteProvider, endpoint, path)
	if err != nil {
		return err
	}
	v.WatchRemoteConfigOnChannel()

	return nil
}
