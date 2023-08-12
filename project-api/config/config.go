package config

import (
	"github.com/spf13/viper"
	"github.com/spxzx/project-common/logs"
	"log"
	"os"
)

type ServerConf struct {
	Name string
	Port string
}

type GRPCConf struct {
	Addr string
}

type EtcdConf struct {
	Addrs []string
}

type Config struct {
	viper  *viper.Viper
	Server *ServerConf
	GRPC   *GRPCConf
	Etcd   *EtcdConf
}

var Conf *Config

func init() {
	conf := &Config{viper: viper.New()}
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
	Conf = conf
}

func (c *Config) initServerConf() {
	c.Server = &ServerConf{
		Name: c.viper.GetString("server.name"),
		Port: c.viper.GetString("server.port"),
	}
}

func (c *Config) initGRPCConf() {
	c.GRPC = &GRPCConf{
		Addr: c.viper.GetString("grpc.addr"),
	}
}
func (c *Config) initEtcdConf() {
	var addrs []string
	c.viper.UnmarshalKey("etcd.addrs", &addrs)
	c.Etcd = &EtcdConf{Addrs: addrs}
}

func (c *Config) initZapLog() {
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
