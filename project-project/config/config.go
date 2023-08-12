package config

import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"github.com/spxzx/project-common/logs"
	"log"
	"os"
	"time"
)

type serverConf struct {
	Name string
	Port string
}

type gRPCConf struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
}

type etcdConf struct {
	Addrs []string
}

type mysqlConf struct {
	Username    string
	Password    string
	Host        string
	Port        int
	Db          string
	TablePrefix string
}

type TokenConf struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}

type dbConf struct {
	Separation bool
	Master     mysqlConf
	Slave      []mysqlConf
}

type config struct {
	viper  *viper.Viper
	Server *serverConf
	GRPC   *gRPCConf
	Etcd   *etcdConf
	Mysql  *mysqlConf
	Token  *TokenConf
	Db     *dbConf
}

var Conf *config

func init() {
	conf := &config{viper: viper.New()}
	nacos := InitNacosClient()
	confYaml, err := nacos.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  BC.NacosConfig.Group,
	})
	if err != nil {
		log.Fatalln(err)
	}
	conf.viper.SetConfigType("yaml")
	log.Println(confYaml, err)
	if confYaml != "" {
		err = conf.viper.ReadConfig(bytes.NewBuffer([]byte(confYaml)))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("load nacos config")
		err = nacos.confClient.ListenConfig(vo.ConfigParam{
			DataId: "config.yaml",
			Group:  BC.NacosConfig.Group,
			OnChange: func(namespace, group, dataId, data string) {
				log.Println("listen nacos config change", data)
				//监听变化
				err = conf.viper.ReadConfig(bytes.NewBuffer([]byte(data)))
				if err != nil {
					log.Printf("listen nacos config parse err %s \n", err.Error())
				}
				//重新载入配置
				conf.ReLoadAllConfig()
			},
		})
	} else {
		workDir, _ := os.Getwd()
		conf.viper.SetConfigName("config")
		//conf.viper.SetConfigType("yaml")
		conf.viper.AddConfigPath("/etc/ProjectManagement/project-grpc")
		conf.viper.AddConfigPath(workDir + "/config")
		if err = conf.viper.ReadInConfig(); err != nil {
			log.Fatalln(err)
		}
	}
	conf.ReLoadAllConfig()
}

func (c *config) ReLoadAllConfig() {
	c.initServerConf()
	c.initGRPCConf()
	c.initEtcdConf()
	c.initZapLog()
	c.initTokenConf()
	c.initDbConfig()
	//重新创建相关的客户端
	c.ReConnRedis()
	c.ReConnMysql()
	Conf = c
}

func (c *config) initServerConf() {
	c.Server = &serverConf{
		Name: c.viper.GetString("server.name"),
		Port: c.viper.GetString("server.port"),
	}
}

func (c *config) initGRPCConf() {
	c.GRPC = &gRPCConf{
		Addr:    c.viper.GetString("grpc.addr"),
		Name:    c.viper.GetString("grpc.name"),
		Version: c.viper.GetString("grpc.version"),
		Weight:  c.viper.GetInt64("grpc.weight"),
	}
}

func (c *config) initEtcdConf() {
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	c.Etcd = &etcdConf{Addrs: addrs}
}

func (c *config) initZapLog() {
	if err := logs.InitLogger(&logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxAge:        c.viper.GetInt("zap.maxAge"),
		MaxSize:       c.viper.GetInt("zap.maxSize"),
		MaxBackups:    c.viper.GetInt("zap.maxBackups"),
	}); err != nil {
		log.Fatalln(err)
	}
}

func (c *config) InitRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

func (c *config) InitMysqlOptions() string {
	c.Mysql = &mysqlConf{TablePrefix: c.viper.GetString("mysql.tablePrefix")}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.viper.GetString("mysql.username"),
		c.viper.GetString("mysql.password"),
		c.viper.GetString("mysql.host"),
		c.viper.GetString("mysql.port"),
		c.viper.GetString("mysql.db"),
	)
}

func (c *config) initTokenConf() {
	c.Token = &TokenConf{
		AccessExp:     time.Duration(c.viper.GetInt("jwt.accessExp")) * 24 * time.Hour,
		RefreshExp:    time.Duration(c.viper.GetInt("jwt.refreshExp")) * 24 * time.Hour,
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
	}
}

func (c *config) initDbConfig() {
	mc := &dbConf{}
	mc.Separation = c.viper.GetBool("db.separation")
	var slaves []mysqlConf
	err := c.viper.UnmarshalKey("db.slave", &slaves)
	if err != nil {
		panic(err)
	}
	master := mysqlConf{
		Username:    c.viper.GetString("db.master.username"),
		Password:    c.viper.GetString("db.master.password"),
		Host:        c.viper.GetString("db.master.host"),
		Port:        c.viper.GetInt("db.master.port"),
		Db:          c.viper.GetString("db.master.db"),
		TablePrefix: c.viper.GetString("db.master.tablePrefix"),
	}
	mc.Master = master
	mc.Slave = slaves
	c.Db = mc
}
