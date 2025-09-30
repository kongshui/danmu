package conf

import (
	"log"
	"os"

	"github.com/kongshui/danmu/common"

	"github.com/pelletier/go-toml/v2"
)

type (
	//Server
	Server struct {
		Addr          string `toml:"addr"`           //节点地址
		Port          string `toml:"port"`           //节点端口
		Name          string `toml:"name"`           //节点名称
		GroupId       string `toml:"group_id"`       //分组id
		NodeType      int8   `toml:"node_type"`      //节点类型
		ListenMode    string `toml:"listen_mode"`    //监听模式
		HeartbeatOpen bool   `toml:"heartbeat_open"` //心跳时间
	}
	//Etcd
	Etcd struct {
		Addr     []string `toml:"addr"`     //etcd地址
		Username string   `toml:"username"` //etcd 用户名
		Password string   `toml:"password"` //etcd 密码
	}
	//Gateway
	Gateway struct {
		ConnMode  string `toml:"mode"`       //链接模式，默认hash
		GroupMode string `toml:"group_mode"` //group modeType
	}
	//Log
	Logging struct {
		Level      string `toml:"level"`
		LogPath    string `toml:"log_path"`
		MaxSize    int64  `toml:"max_size"`
		MaxBackups int    `toml:"max_backups"`
		MaxAge     int    `toml:"max_age"`
	}
	//Conf
	Config struct {
		Project string  `toml:"project"`
		Server  Server  `toml:"server"`
		Logging Logging `toml:"logging"`
		Etcd    Etcd    `toml:"etcd"`
		Gateway Gateway `toml:"gateway"`
	}
)

// conf new
func newConf() *Config {
	return &Config{
		Server:  Server{"127.0.0.1", "6666", "gateway1", "0", 1, "websocket", false},
		Logging: Logging{"", "", 0, 0, 0},
		Etcd:    Etcd{[]string{"127.0.0.1:2379"}, "root", "123456"},
		Gateway: Gateway{"hash", ""},
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
