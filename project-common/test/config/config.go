package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
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
	TablePrefix string
}

type TokenConf struct {
	AccessExp     time.Duration
	RefreshExp    time.Duration
	AccessSecret  string
	RefreshSecret string
}

type config struct {
	viper  *viper.Viper
	Server *serverConf
	GRPC   *gRPCConf
	Etcd   *etcdConf
	Mysql  *mysqlConf
	Token  *TokenConf
}

var Conf *config

func init() {
	conf := &config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("/etc/ProjectManagement/project-grpc")
	conf.viper.AddConfigPath(workDir + "/config")

	if err := conf.viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	conf.initServerConf()
	conf.initGRPCConf()
	conf.initEtcdConf()
	conf.initZapLog()
	conf.initTokenConf()
	Conf = conf
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
