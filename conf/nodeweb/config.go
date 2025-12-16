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
		Addr       string `toml:"addr"`        //节点地址
		Port       string `toml:"port"`        //节点端口
		Name       string `toml:"name"`        //节点名称
		GroupId    string `toml:"group_id"`    //分组id
		NodeType   int8   `toml:"node_type"`   //节点类型
		ListenMode string `toml:"listen_mode"` //监听模式
		WebPort    string `toml:"web_port"`    //web 端口
	}
	// WEb
	Web struct {
		Addr  string `toml:"addr"`  //节点地址
		Port  string `toml:"port"`  //节点端口
		Debug bool   `toml:"debug"` //是否调试
		Grpc  bool   `toml:"grpc"`  //是否开启grpc
	}
	// App 应用配置
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
		ConfigDir         string `toml:"config_dir"`          // 配置文件目录
		LogStoreDir       string `toml:"log_store_dir"`       // 存储日志目录
		UserChangeTime    int64  `toml:"user_change_time"`    // 用户信息变更时间,单位秒
		DeleteDay         int    `toml:"delete_day"`          // 删除多少天之前的日志文件
		IsAnonymous       bool   `toml:"is_anonymous"`        // 是否匿名
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
	// Log
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
		IsNode  bool    `toml:"is_node"`
		Server  Server  `toml:"server"`
		Web     Web     `toml:"web"`
		Logging Logging `toml:"logging"`
		Etcd    Etcd    `toml:"etcd"`
		Redis   Redis   `toml:"redis"`
		Nats    Nats    `toml:"nats"`
		Mysql   Mysql   `toml:"mysql"`
		App     App     `toml:"app"` //应用配置
	}
)

// GetConf implements conf.Config.

// conf new
func newConf() *Config {
	return &Config{
		Server:  Server{"127.0.0.1", "5555", "node1", "0", 2, "tcp", "8080"},
		Web:     Web{"127.0.0.1", "8080", false, false},
		Logging: Logging{"", "", 0, 0, 0},
		Etcd:    Etcd{[]string{"127.0.0.1:2379"}, "root", "123456"},
		Nats:    Nats{[]string{"nats://localhost:4222"}},
		App:     App{"ks", "", "", "", "", "", "", "", "", "", "./logstore", 300, 7, true, false, false, false},
		Redis:   Redis{"127.0.0.1:6379", 0, "", false},
		Mysql:   Mysql{"127.0.0.1:3306", "root", "123456", "store_log", true},
		IsNode:  true,
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
