package orm

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"strings"
	"time"
)

type dbConfig struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	DbName   string `json:"dbname"`
}

type groups struct {
	Name string
	Key  string
}

type emptyVal struct {
	Name  string
	Value string
	File  string
}

func init() {

	databaseConf, err := getConf("database")
	if err != nil {
		panic(err)
	}
	appConf, err := getConf("app")
	if err != nil {
		panic(err)
	}

	mysqlKey := databaseConf.String("mysqlkey")
	appName := databaseConf.String("appname")
	appGroup := databaseConf.String("appgroup")
	if len(appName) == 0 {
		appName = appConf.String("appname")
	}
	if len(appGroup) == 0 {
		appGroup = appConf.String("appgroup")
	}
	//db
	DbName := databaseConf.String("DbName")
	UserName := databaseConf.String("UserName")
	Host := databaseConf.String("Host")
	Port, _ := databaseConf.Int64("Port")
	Password := databaseConf.String("Password")

	checkEmpty(emptyVal{"appname", appName, "app.conf"}, emptyVal{"appgroup", appGroup, "app.conf"})

	orm.RegisterDriver("mysql", orm.DRMySQL)

	var dbGroup map[string]groups
	dbGroup = make(map[string]groups)

	if len(mysqlKey) != 0 {
		// 默认设置
		dbGroup["default"] = groups{"default", mysqlKey}
	}

	// 支持多组
	//addMultiple(databaseConf, dbGroup)

	/*获取多个数据库配置*/
	//configs := getDbConfigs(dbGroup, appGroup, appName)
	configs := getDbConfigsOne(dbGroup,DbName,UserName,Host,Port,Password)

	//registerDev(configs, databaseConf)
	/*
		if len(configs) == 0 {
			panic("there are no databases available")
		}



		if !checkDefaultDb(configs) {
			panic("must contain default database")
		}
	*/
	//统一注册
	register(configs)
}

func addMultiple(databaseConf config.Configer, dbGroup map[string]groups) {
	mysqlKeyMultiple := databaseConf.String("mysqlkeymultiple")

	if len(mysqlKeyMultiple) > 0 {
		multiple := strings.Split(mysqlKeyMultiple, ",")

		for _, atom := range multiple {

			group := strings.Split(atom, "::")
			if len(group) != 2 || len(group[0]) == 0 || len(group[1]) == 0 {
				panic("database.conf mysqlkeymultiple format error,for example:  mysqlkeymultiple = name1::key1,name2::key2")
			}

			dbGroup[group[0]] = groups{group[0], group[1]}
		}
	}
}

func checkDefaultDb(configs map[string]dbConfig) bool {
	for name, _ := range configs {
		if name == "default" {
			return true
		}
	}
	return false
}
func getDbConfigsOne(dbGroup map[string]groups,DbName string ,UserName string ,Host string ,Port int64 ,Password string) map[string]dbConfig {
	var configs map[string]dbConfig
	configs = make(map[string]dbConfig)
	var dbconfig dbConfig
	dbconfig.DbName= DbName
	dbconfig.UserName= UserName
	dbconfig.Host= Host
	dbconfig.Port= Port
	//err := json.Unmarshal([]byte(content), &dbconfig)
	//if err != nil {
	//	panic(err)
	//}
	//var dbconfig dbConfig
	dbconfig.Password= Password
	configs["default"] = dbconfig


	return configs
}
func getDbConfigs(dbGroup map[string]groups, appGroup string, appName string) map[string]dbConfig {
	var configs map[string]dbConfig
	configs = make(map[string]dbConfig)

	for _, group := range dbGroup {

		var dbconfig dbConfig
		dbconfig.DbName="test"
		dbconfig.UserName="root"
		dbconfig.Host="localhost"
		dbconfig.Port=3306
		//err := json.Unmarshal([]byte(content), &dbconfig)
		//if err != nil {
		//	panic(err)
		//}
		//var dbconfig dbConfig
		dbconfig.Password="123456"
		configs[group.Name] = dbconfig
	}

	return configs
}

func register(configs map[string]dbConfig) {
	for name, dbconfig := range configs {
		orm.RegisterDataBase(name, "mysql", dbconfig.UserName+":"+dbconfig.Password+"@tcp("+dbconfig.Host+":"+strconv.FormatInt(dbconfig.Port, 10)+")/"+dbconfig.DbName)

		//设置连接超时为25s，比目前测试、生产的30s小即可，防止连接池可用连接不可用
		db,_ := orm.GetDB(name)
		db.SetConnMaxLifetime(25*time.Second)
	}
}

func registerDev(configs map[string]dbConfig, databaseConf config.Configer) {

	if !isProduct() {
		mysqlhost := databaseConf.String("mysqlhost")
		mysqluser := databaseConf.String("mysqluser")
		mysqlpass := databaseConf.String("mysqlpass")
		mysqldb := databaseConf.String("mysqldb")

		mysqlport, _ := databaseConf.Int64("mysqlport")
		if mysqlport == 0 {
			mysqlport = 3306
		}

		// 使用本地开发模式
		if len(mysqlhost) > 0 || len(mysqluser) > 0 || len(mysqlpass) > 0 || len(mysqldb) > 0 {

			checkEmpty(
				emptyVal{"mysqlhost", mysqlhost, "database.conf"},
				emptyVal{"mysqluser", mysqluser, "database.conf"},
				emptyVal{"mysqlpass", mysqlpass, "database.conf"},
				emptyVal{"mysqldb", mysqldb, "database.conf"},
			)

			//覆盖default
			configs["default"] = dbConfig{
				UserName: mysqluser,
				Password: mysqlpass,
				Host:     mysqlhost,
				Port:     mysqlport,
				DbName:   mysqldb,
			}
		}
	}
}

func getConf(name string) (config.Configer, error) {
	projectRoot := os.Getenv("PROJECT_ROOT")
	confDir := ""
	if len(projectRoot) == 0 {
		confDir = "conf/" + name + ".conf"
	} else {
		confDir = projectRoot + "conf/" + name + ".conf"
	}
	databaseConf, err := config.NewConfig("ini", confDir)
	return databaseConf, err
}

func checkEmpty(vals ...emptyVal) {
	//检查空值
	for _, val := range vals {
		if len(val.Value) == 0 {
			panic(fmt.Errorf("in the file %s, %s cannot be empty", val.File, val.Name))
		}
	}
}

func isProduct() bool {
	env := os.Getenv("K8S_CLUSTER_TYPE")
	if env == "product" {
		return true
	}
	return false
}
