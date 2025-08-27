package dao

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	Client *sql.DB
	Lock   *sync.RWMutex
	isUse  bool
}

// new mysql client
func NewMysqlClient() *MysqlClient {
	return &MysqlClient{
		Client: nil,
		Lock:   &sync.RWMutex{},
	}
}

func init() {
}

// mysql初始化
func (m *MysqlClient) MysqlInit(userName, password, addr, db string, isUse bool) {
	if !isUse {
		return
	}
	if addr == "" {
		log.Println("从环境中获取MYSQL地址失败...")
		os.Exit(6)
	}
	dsn := userName + ":" + password + "@tcp(" + addr + ")/" + db + "?charset=utf8mb4&parseTime=True"
	client, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	m.Client = client
	m.isUse = isUse
	go m.mysqlCheckPing(userName, password, addr, db)
}

// mysql连接检查
func (m *MysqlClient) mysqlCheckPing(userName, password, addr, db string) {
	if !m.isUse {
		return
	}
	t := time.NewTicker(10 * time.Second)
	for {
		<-t.C
		if err := m.Client.Ping(); err != nil {
			log.Println("mysql连接失败,重新初始化")
			m.MysqlInit(userName, password, addr, db, m.isUse)
			continue
		}
		// log.Println("mysql ping success")
	}
}

// mysql插入礼物数据roomid, anchorOpenId, open_id, nick_name, msg_id string, giftNum, giftValue int, isTest bool, gift_time int64
func (m *MysqlClient) InsertGiftData(roomid, anchorOpenId, anchorName, roundId, userId, userName, msgId, giftId string, giftCount, giftValue int, isTest bool) error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("INSERT INTO log_gift (room_id, anchor_open_id,anchor_name,open_id, nick_name, msg_id, gift_num,gift_value,test,gift_time,gift_id,round_id) VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		roomid, anchorOpenId, anchorName, userId, userName, msgId, giftCount, giftValue, isTest, time.Now().UnixMilli(), giftId, roundId)
	if err != nil {
		log.Println("mysql插入礼物数据失败", err)
	}
	return err
}

// mysql 查询维护状态
func (m *MysqlClient) GetMaintainStatus() (WeihuStruct, error) {
	if !m.isUse {
		return WeihuStruct{}, nil
	}
	var status WeihuStruct
	err := m.Client.QueryRow("SELECT is_maintain,start_time,end_time,maintain_msg FROM maintain").Scan(&status.IsMaintain, &status.StartTime, &status.EndTime, &status.MaintainMsg)
	if err != nil {
		log.Println("mysql查询维护状态失败", err)
		return WeihuStruct{}, err
	}
	return status, nil
}

// mysql 当前排行版清零
func (m *MysqlClient) ClearRank() error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("UPDATE player_info SET world_point = 0")
	if err != nil {
		log.Println("mysql 当前排行版清零失败", err)
		return errors.New("mysql 当前排行版清零失败: " + err.Error())
	}
	return nil
}

// mysql 当前连胜币清零
func (m *MysqlClient) ClearCoin() error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("UPDATE player_info SET win_coin = 0")
	if err != nil {
		log.Println("mysql 当前连胜币清零失败", err)
		return errors.New("mysql 当前连胜币清零失败: " + err.Error())
	}
	return nil
}

// mysql update 排行榜
func (m *MysqlClient) UpdateRank(openid string, score int64) error {
	if !m.isUse {
		return nil
	}
	if int64(score) < 1 {
		return nil
	}
	// result, err := m.Client.Exec("UPDATE player_info SET world_point = world_point+?,world_history_point = world_history_point+? WHERE open_id = ?", score, score, openid)
	result, err := m.Client.Exec("insert into player_info (open_id,world_point,world_history_point) values(?,?,?) on duplicate key update world_point = world_point+?,world_history_point = world_history_point+?", openid, score, score, score, score)
	if err != nil {
		log.Println("mysql update 排行榜失败", err)
		return errors.New("mysql update 排行榜失败: " + err.Error())
	}
	statInt, err := result.RowsAffected()
	if statInt == 0 || err != nil {
		return errors.New("mysql update 排行榜失败: 没有找到openid或者没有更新， openID： " + openid + ", score： " + strconv.FormatInt(score, 10))
	}
	return nil
}

// mysql 设置玩家连胜币
func (m *MysqlClient) SetCoin(openid string, score int64) error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("UPDATE player_info SET win_coin = ? WHERE open_id = ?", score, openid)
	if err != nil {
		log.Println("mysql 设置玩家连胜币失败", err)
		return errors.New("mysql 设置玩家连胜币失败: " + err.Error())
	}
	return nil
}

// mysql 更新玩家连胜币
func (m *MysqlClient) UpdateCoin(openid string, score int64) error {
	if !m.isUse {
		return nil
	}
	if int64(score) == 0 {
		return nil
	}
	// ok, err := m.IsPlayerExist(openid)
	// if err != nil {
	// 	log.Println("mysql update 连胜币失败", err)
	// 	return err
	// }
	// if !ok {
	// 	_, err := m.Client.Exec("insert into player_info (open_id,win_coin) values(?,?)", openid, score)
	// 	if err != nil {
	// 		if !strings.Contains(err.Error(), "1062") {
	// 			log.Println("mysql 插入玩家信息失败,", err, openid, score)
	// 			return errors.New("mysql update 连胜币失败: " + err.Error())
	// 		}
	// 	}
	// }

	result, err := m.Client.Exec("insert into player_info (open_id,win_coin) values(?,?) on duplicate key update win_coin = win_coin+?", openid, score, score)
	// result, err := m.Client.Exec("UPDATE player_info SET win_coin = win_coin+? WHERE open_id = ?", score, openid)
	if err != nil {
		log.Println("mysql update 连胜币失败", err)
		return errors.New("mysql update 连胜币失败: " + err.Error())
	}
	statInt, err := result.RowsAffected()
	if statInt == 0 || err != nil {
		return errors.New("mysql update 连胜币失败: 没有找到openid或者没有更新， openID： " + openid + ", score： " + strconv.FormatInt(score, 10))
	}
	return nil
}

// mysql 更新玩家连胜
func (m *MysqlClient) UpdateWin(openid string, score int64) error {
	if !m.isUse {
		return nil
	}
	if int64(score) == 0 {
		return nil
	}
	result, err := m.Client.Exec("UPDATE player_info SET win_count = win_count+? WHERE open_id = ?", score, openid)
	if err != nil {
		log.Println("mysql update 连胜失败", err)
		return errors.New("mysql update 连胜失败: " + err.Error())
	}
	statInt, err := result.RowsAffected()
	if statInt == 0 || err != nil {
		return errors.New("mysql update 连胜失败: 没有找到openid或者没有更新， openID： " + openid + ", score： " + strconv.FormatInt(score, 10))
	}
	return nil
}

// mysql 设置玩家连胜
func (m *MysqlClient) SetWin(openid string, score int64) error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("UPDATE player_info SET win_count = ? WHERE open_id = ?", score, openid)
	if err != nil {
		log.Println("mysql 设置玩家连胜失败", err)
		return errors.New("mysql 设置玩家连胜失败: " + err.Error())
	}
	return nil
}

// 查询玩家是否存在
func (m *MysqlClient) IsPlayerExist(openid string) (bool, error) {
	if !m.isUse {
		return false, nil
	}
	var count int
	err := m.Client.QueryRow("SELECT COUNT(1) FROM player_info WHERE open_id = ?", openid).Scan(&count)
	if err != nil {
		log.Println("mysql查询玩家是否存在失败", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// 插入玩家基本信息
func (m *MysqlClient) InsertPlayerBaseInfo(openId, avatarUrl, nickName string) error {
	if !m.isUse {
		return nil
	}
	if ok, _ := m.IsPlayerExist(openId); ok {
		return nil
	}
	_, err := m.Client.Exec("insert into player_info (open_id,avatar_url,nick_name) values(?,?,?)", openId, avatarUrl, nickName)
	if err != nil {
		if !strings.Contains(err.Error(), "1062") {
			log.Println("mysql 插入玩家信息失败", err)
			return errors.New("mysql 插入玩家信息失败: " + err.Error())
		}
	}
	return nil
}

// 更新玩家基本信息
func (m *MysqlClient) UpdatePlayerBaseInfo(openId, avatarUrl, nickName string) error {
	if !m.isUse {
		return nil
	}
	_, err := m.Client.Exec("update player_info set avatar_url = ?,nick_name = ? where open_id = ?", avatarUrl, nickName, openId)
	if err != nil {
		log.Println("mysql 更新玩家信息失败", err)
		return errors.New("mysql 更新玩家信息失败: " + err.Error())
	}
	return nil
}

// 维护信息
type WeihuStruct struct {
	IsMaintain  bool   // 是否维护中
	StartTime   string // 开始时间
	EndTime     string // 结束时间
	MaintainMsg string // 维护信息
}

// 添加玩家获胜统计,about： true 左边，false，右边
func (m *MysqlClient) UpdateOpenWinCount(openid string, about bool) error {
	if !m.isUse {
		return nil
	}
	// ok, err := m.isInWinCount(openid)
	// if err != nil {
	// 	return err
	// }
	var context string
	if about {
		context = "left_count"
	} else {
		context = "right_count"
	}
	_, err := m.Client.Exec("insert into winner_count (open_id,"+context+") values(?,?) on duplicate key update "+context+" = "+context+"+1", openid, 1)
	if err != nil {
		log.Println("mysql 添加玩家获胜统计失败", err)
		return errors.New("mysql 添加玩家获胜统计失败: " + err.Error())
	}
	return nil
}

// 查询玩家是否在获胜统计中
func (m *MysqlClient) IsInWinCount(openid string) (bool, error) {
	if !m.isUse {
		return false, nil
	}
	var count int
	err := m.Client.QueryRow("SELECT COUNT(1) FROM winner_count WHERE open_id = ?", openid).Scan(&count)
	if err != nil {
		log.Println("mysql查询玩家是否在获胜统计中失败", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// 更新左右胜场数， about: true 左边，false，右边
func (m *MysqlClient) UpdateGroupWinCount(groupId string) error {
	if !m.isUse {
		return nil
	}

	_, err := m.Client.Exec("INSERT INTO win_group_count (group_id,win_count) VALUES (?,1) ON DUPLICATE KEY UPDATE win_count = win_count+1", groupId)
	if err != nil {
		log.Println("mysql 更新左右胜场数失败", err)
		return errors.New("mysql 更新左右胜场数失败: " + err.Error())
	}
	return nil
}

// 查询组的胜场次数
func (m *MysqlClient) QueryGroupWinCount(groupId string) (int64, error) {
	if !m.isUse {
		return 0, nil
	}
	var count int64
	err := m.Client.QueryRow("SELECT win_count FROM win_group_count WHERE group_id = ?", groupId).Scan(&count)
	if err != nil {
		log.Println("mysql 查询组的胜场次数失败", err)
		return 0, err
	}
	return count, nil
}

// 查询玩家信息
func (m *MysqlClient) QueryPlayerInfo(openid string) (string, string, error) {
	if !m.isUse {
		return "", "", nil
	}
	var (
		avatarUrl string
		nickName  string
	)
	err := m.Client.QueryRow("SELECT avatar_url,nick_name FROM player_info WHERE open_id = ?", openid).Scan(&avatarUrl, &nickName)
	if err != nil {
		log.Println("mysql 查询玩家信息失败", err)
		return "", "", err
	}
	return avatarUrl, nickName, nil
}
