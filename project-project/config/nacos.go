package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

type NacosClient struct {
	confClient config_client.IConfigClient
}

func InitNacosClient() *NacosClient {
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				NamespaceId:         BC.NacosConfig.Namespace,
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "/tmp/nacos/log",
				CacheDir:            "/tmp/nacos/cache",
				LogLevel:            "debug",
			},
			ServerConfigs: []constant.ServerConfig{
				{
					IpAddr:      BC.NacosConfig.IpAddr,
					ContextPath: BC.NacosConfig.ContextPath,
					Port:        uint64(BC.NacosConfig.Port),
					Scheme:      BC.NacosConfig.Scheme,
				},
			},
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	nc := &NacosClient{
		confClient: configClient,
	}
	return nc
}
