# qsmysql

1. If you use gorm for mysql connection and use viper as major configuration parse tool, you may like my project

2. First, you mysql configurtion may like this below:

   ```yaml
   mysql:
     master:
       host: localhost
       port: 3306
       user: root
       password:
       db: test
       charset: utf8mb4
       max_idle_conns: 2
       max_open_conns: 5
       log_mode: true
       ## minutes
       conn_max_lifetime: 60
     slave:
       host:
         - localhost
         - localhost1
       port: 3306
       user: root
       password:
       db: test
       charset: utf8mb4
       max_idle_conns: 2
       max_open_conns: 5
       log_mode: true
       ## minutes
       conn_max_lifetime: 60
   ```

   this section of mysql configuration can be any part of your configuration file.

3. Then, you can just use this package as:

   ```go
   package main
   
   import (
   	"reflect"
   
   	"github.com/chunshengster/qsmysql"
   	"github.com/spf13/viper"
   )
   
   func main() {
   	viper.SetConfigFile("./config.yaml")
   	viper.ReadInConfig()
   
   	qsmysql.SetConfig(viper.Sub("mysql"))
     defer qsmysql.Close()
   	master := qsmysql.GetMaster()  // instance of *gorm.DB
   	reflect.TypeOf(master)
   }
   ```

   

