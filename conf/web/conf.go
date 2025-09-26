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
		Addr  string `toml:"addr"`  //节点地址
		Port  string `toml:"port"`  //节点端口
		Debug bool   `toml:"debug"` //是否调试
		Grpc  bool   `toml:"grpc"`  //是否开启grpc
	}
	App struct {
		PlatForm          string `toml:"platform"`            //平台
		AppId             string `toml:"app_id"`              //appId
		AppSecret         string `toml:"app_secret"`          //appSecret
		CommentSecret     string `toml:"comment_secret"`      //评论密钥
		GiftSecret        string `toml:"gift_secret"`         //礼物密钥
		LikeSecret        string `toml:"like_secret"`         //点赞密钥
		ChooseGroupSecret string `toml:"choose_group_secret"` //选择组密钥
		QueryGroupSecret  string `toml:"query_group_secret"`  //查询组密钥
		RaceSecret        string `toml:"race_secret"`         // 比赛密钥
		IsOnline          bool   `toml:"is_online"`           // 是否上线
		IsMock            bool   `toml:"is_mock"`             // 是否模拟
		NoSend            bool   `toml:"no_send"`             // 是否不发送
	}
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
	//Mysql
	Mysql struct {
		Addr     string `toml:"addr"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Db       string `toml:"db"`
		IsUse    bool   `toml:"isuse"`
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
		Project string  `toml:"project"`
		Server  Server  `toml:"server"`
		Logging Logging `toml:"logging"`
		Etcd    Etcd    `toml:"etcd"`
		Redis   Redis   `toml:"redis"`
		Nats    Nats    `toml:"nats"`
		Mysql   Mysql   `toml:"mysql"`
		App     App     `toml:"app"` //应用配置
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
		Server:  Server{"127.0.0.1", "5555", true, true},
		Logging: Logging{"", "", 0, 0, 0},
		Etcd:    Etcd{[]string{"127.0.0.1:2379"}, "root", "123456"},
		Nats:    Nats{[]string{"nats://localhost:4222"}},
		Redis:   Redis{"127.0.0.1:6379", 0, "", false},
		Mysql:   Mysql{"127.0.0.1:3306", "root", "123456", "store_log", true},
		App:     App{"ks", "", "", "", "", "", "", "", "", false, false, false},
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
