package node

import (
	"encoding/json"
	"fmt"
	"gaad/common"
	"gaad/controllers"
	"gaad/db/sqlitedb"
	"gaad/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"strconv"
)

// 查看全部在线用户
func CreateNode(c *gin.Context) {

	data, _ := ioutil.ReadAll(c.Request.Body)

	node := &models.Node{}

	err := json.Unmarshal(data, node)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}

	address := net.ParseIP(node.Ip)
	if address == nil {
		controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
		return
	} else {
		fmt.Println("正确的ip地址", address.String())
	}

	if node.Port < 2 {
		controllers.Response(c, common.ParameterIllegal, "端口不可以小于2", nil)
		return
	}

	if node.Username == "" || node.Password == "" {
		controllers.Response(c, common.ParameterIllegal, "", nil)
		return
	}
	sqlitedb.Create(&node)

	controllers.Response(c, common.OK, "", nil)

}

// 查看全部在线用户
func ListNodes(c *gin.Context) {

	curPageStr := c.Query("curPage")
	pageRecordStr := c.Query("pageRecord")

	if curPageStr == "" {
		curPageStr = "1"
	}
	if pageRecordStr == "" {
		pageRecordStr = "10"
	}

	curPage, _ := strconv.Atoi(curPageStr)
	pageRecord, _ := strconv.Atoi(pageRecordStr)

	var (
		nodes []models.Node
	)

	total := sqlitedb.QueryPage(curPage, pageRecord, &nodes, "username = ?", "root")

	data := make(map[string]interface{})

	data["nodes"] = nodes
	data["curPage"] = curPage
	data["pageRecord"] = pageRecord
	data["total"] = total
	controllers.Response(c, common.OK, "", data)
}
