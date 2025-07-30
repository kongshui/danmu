package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// 服务维护
func MaintainHandle(c *gin.Context) {
	data, _ := mysql.GetMaintainStatus()
	d, _ := json.Marshal(data)
	var sData any
	json.Unmarshal(d, &sData)
	c.JSON(200, sData)
}
