package common

type (
	RoomRegister struct {
		Uuid       string // client uuid
		OpenId     string //OpenId
		UserId     string //相当于OpenId
		RoomId     string // roomid
		GroupId    string // groupid,组Id
		GradeLevel int64  // 等级
	}
)
