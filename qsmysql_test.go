package qsmysql

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var conf []*viper.Viper

// func TestMain(t *testing.M) {

// }

func TestQSMySQL_setConfig(t *testing.T) {
	config_files := []string{
		"./config.test1.yaml",
		"./config.test2.yaml",
	}
	for _, c := range config_files {
		viper.SetConfigFile(c)
		viper.ReadInConfig()
		x := viper.Sub("mysql")
		conf = append(conf, x)
		viper.Reset()
	}
	type fields struct {
		master     *mysqlMaster
		slave      *mysqlSlave
		viper_conf *viper.Viper
	}
	type args struct {
		v *viper.Viper
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "one master one slave",
			fields: fields{
				master: &mysqlMaster{
					handler: nil,
					dsn:     "test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local",
					conf: connPoolConf{
						max_open_conns:    5,
						max_idle_conns:    2,
						conn_max_lifetime: 60,
						log_mode:          true,
					},
				},
				slave: &mysqlSlave{
					handlers: nil,
					dsns:     []string{"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local"},
					conf: connPoolConf{
						max_open_conns:    5,
						max_idle_conns:    2,
						conn_max_lifetime: 60,
						log_mode:          true,
					},
				},
				viper_conf: conf[0],
			},
			args:    args{v: conf[0]},
			wantErr: false,
		},
		{
			name: "one master two slave",
			fields: fields{
				master: &mysqlMaster{
					handler: nil,
					dsn:     "test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local",
					conf: connPoolConf{
						max_open_conns:    5,
						max_idle_conns:    2,
						conn_max_lifetime: 60,
						log_mode:          true,
					},
				},
				slave: &mysqlSlave{
					handlers: nil,
					dsns: []string{"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local",
						"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wb.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local"},
					conf: connPoolConf{
						max_open_conns:    5,
						max_idle_conns:    2,
						conn_max_lifetime: 60,
						log_mode:          true,
					},
				},
				viper_conf: conf[1],
			},
			args:    args{v: conf[1]},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QSMySQL{
				master:     tt.fields.master,
				slave:      tt.fields.slave,
				viper_conf: tt.fields.viper_conf,
			}
			t.Logf("slave.slave.dsns = %v", q.slave.dsns)
			if err := q.SetConfig(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("QSMySQL.setConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(q.master, tt.fields.master) {
				t.Errorf("q.master = %v, want %v", q.master, tt.fields.master)
			}
			if !reflect.DeepEqual(q.slave.handlers, tt.fields.slave.handlers) {
				t.Errorf("q.slave.handlers = %v, want %v", q.slave.handlers, tt.fields.slave.handlers)
			}
			if !reflect.DeepEqual(q.slave.dsns, tt.fields.slave.dsns) {
				t.Errorf("q.slave.dsns = %v, want %v", q.slave.dsns, tt.fields.slave.dsns)
			}
			if !reflect.DeepEqual(q.slave.conf, tt.fields.slave.conf) {
				t.Errorf("q.slave.conf = %v,want %v", q.slave.conf, tt.fields.slave.conf)
			}
			if !reflect.DeepEqual(q.viper_conf, tt.fields.viper_conf) {
				t.Errorf("q.viper_conf = %v, want %v", q.viper_conf, tt.fields.viper_conf)
			}

		})
	}
}

func TestSetConfig(t *testing.T) {
	config_files := []string{
		"./config.test1.yaml",
		"./config.test2.yaml",
	}
	for _, c := range config_files {
		viper.SetConfigFile(c)
		viper.ReadInConfig()
		x := viper.Sub("mysql")
		conf = append(conf, x)
		viper.Reset()
	}
	type args struct {
		v *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "one master one slave config",
			args:    args{v: conf[0]},
			wantErr: false,
		}, {
			name:    "one master two slave config",
			args:    args{v: conf[1]},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetConfig(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseViper(t *testing.T) {
	config_files := []string{
		"./config.test1.yaml",
		"./config.test2.yaml",
	}
	for _, c := range config_files {
		viper.SetConfigFile(c)
		viper.ReadInConfig()
		x := viper.Sub("mysql")
		conf = append(conf, x)
		viper.Reset()
	}
	type args struct {
		v *viper.Viper
	}
	tests := []struct {
		name     string
		args     args
		wantDsns []string
		wantConf connPoolConf
		wantErr  bool
	}{
		{
			name:     "one master",
			args:     args{v: conf[0].Sub("master")},
			wantDsns: []string{"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local"},
			wantConf: connPoolConf{
				max_open_conns:    5,
				max_idle_conns:    2,
				conn_max_lifetime: 60,
				log_mode:          true,
			},
			wantErr: false,
		},
		{
			name:     "one slave",
			args:     args{v: conf[0].Sub("slave")},
			wantDsns: []string{"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local"},
			wantConf: connPoolConf{
				max_open_conns:    5,
				max_idle_conns:    2,
				conn_max_lifetime: 60,
				log_mode:          true,
			},
			wantErr: false,
		},
		{
			name: "two slave",
			args: args{v: conf[1].Sub("slave")},
			wantDsns: []string{"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local",
				"test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wb.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local"},
			wantConf: connPoolConf{
				max_open_conns:    5,
				max_idle_conns:    2,
				conn_max_lifetime: 60,
				log_mode:          true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("tt.args.v = %v", tt.args.v.Get("host"))
			gotDsns, gotConf, err := parseViper(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseViper() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotDsns, tt.wantDsns) {
				t.Errorf("parseViper() gotDsns = %v, want %v", gotDsns, tt.wantDsns)
			}
			if !reflect.DeepEqual(gotConf, tt.wantConf) {
				t.Errorf("parseViper() gotConf = %v, want %v", gotConf, tt.wantConf)
			}
		})
	}
}

func Test_connDB(t *testing.T) {
	type args struct {
		dsn  string
		conf connPoolConf
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.DB
		wantErr bool
	}{
		{
			name: "one server",
			args: args{
				dsn: "test_order:61NNT9RJSLwGelGy@tcp(rm-bp1y256043o82d4wa.mysql.rds.aliyuncs.com:3306)/mutual_insure2?charset=utf8mb4&parseTime=True&loc=Local",
				conf: connPoolConf{
					max_open_conns:    5,
					max_idle_conns:    2,
					conn_max_lifetime: 60,
					log_mode:          true,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnDB(tt.args.dsn, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("connDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got.DB().Stats().MaxOpenConnections == tt.args.conf.max_open_conns)
			if got.DB().Stats().MaxOpenConnections != tt.args.conf.max_open_conns {
				t.Errorf("connPoolConf max_open_conns error = %v, want %v", got.DB().Stats().MaxOpenConnections, tt.args.conf.max_open_conns)
			}
		})
	}
}

func TestGetMaster(t *testing.T) {
	config_files := []string{
		"./config.test1.yaml",
		// "./config.test2.yaml",
	}
	for _, c := range config_files {
		viper.SetConfigFile(c)
		viper.ReadInConfig()
		x := viper.Sub("mysql")
		conf = append(conf, x)
		viper.Reset()
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "get master database",
			want: reflect.TypeOf(&gorm.DB{}).String(),
		},
	}
	for _, tt := range tests {
		SetConfig(conf[0])
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaster(); reflect.TypeOf(got).String() != tt.want {
				t.Errorf("GetMaster() = %v, want %v", got, tt.want)
			} else {
				t.Logf("GetMaster got %v", got)
			}
		})
	}
}

func TestGetSlave(t *testing.T) {
	tests := []struct {
		name string
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSlave(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlave() = %v, want %v", got, tt.want)
			}
		})
	}
}
