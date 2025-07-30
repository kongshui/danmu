package conf

import (
	"log"
	"os"

	"github.com/kongshui/danmu/common"

	"github.com/pelletier/go-toml/v2"
)

type (
	//Etcd
	Etcd struct {
		Addr     []string `toml:"addr"`     //etcd地址
		Username string   `toml:"username"` //etcd 用户名
		Password string   `toml:"password"` //etcd 密码
	}
	//Redis
	Redis struct {
		Addr      string `toml:"addr"`
		Db        int    `toml:"db"`
		Password  string `toml:"password"`
		IsCluster bool   `toml:"mode"`
	}
	//Log
	Logging struct {
		Level      string `toml:"level"`
		LogPath    string `toml:"log_path"`
		MaxSize    int64  `toml:"max_size"`
		MaxBackups int    `toml:"max_backups"`
		MaxAge     int    `toml:"max_age"`
	}
	//nats
	Nats struct {
		Addr []string `toml:"addr"` //nats 地址
	}
	//Conf
	Config struct {
		Logging Logging `toml:"logging"`
		Etcd    Etcd    `toml:"etcd"`
		Redis   Redis   `toml:"redis"`
		Nats    Nats    `toml:"nats"`
	}
	//Gateway
	Gateway struct {
		ConnMode  string `toml:"mode"`       //链接模式，默认hash
		GroupMode string `toml:"group_mode"` //group modeType
	}
)

// conf new
func newConf() *Config {
	return &Config{
		Logging: Logging{"", "", 0, 0, 0},
		Etcd:    Etcd{[]string{"127.0.0.1:2379"}, "root", "123456"},
		Nats:    Nats{[]string{"nats://localhost:4222"}},
		Redis:   Redis{"127.0.0.1:6379", 0, "", false},
	}
}

// conf 读取yaml 文件
func (c *Config) readConf() {
	if !common.PathExists("./config.toml") {
		return
	}
	configFile, err := os.ReadFile("./config.toml")
	if err != nil {
		log.Fatal(err)
	}
	if err := toml.Unmarshal(configFile, c); err != nil {
		log.Fatal(err)
	}
}

// 重新加载配置
func (c *Config) ReloadConf() {
	c.readConf()
}

// 获取配置
func GetConf() *Config {
	conf := newConf()
	conf.readConf()
	return conf
}
