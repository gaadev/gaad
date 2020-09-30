package node

import (
	"fmt"
	"gaad/common"
	"gaad/common/web"
	"gaad/controllers"
	"gaad/db/sqlitedb"
	"gaad/models"
	"github.com/gin-gonic/gin"
	"net"
	"strconv"
)

// 查看全部在线用户
func CreateNode(c *gin.Context) {

	modNode("create", c)

}

// 查看全部在线用户
func UpdateNode(c *gin.Context) {

	modNode("update", c)

}

// 查看全部在线用户
func DeleteNode(c *gin.Context) {

	node := models.Node{}
	web.GetModel(&node, c)

	sqlitedb.Delete(&node)
	controllers.Response(c, common.OK, "", nil)
}

func modNode(operation string, c *gin.Context) {
	node := models.Node{}
	web.GetModel(&node, c)

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
	if operation == "create" {
		sqlitedb.Create(&node)
	}
	if operation == "update" {
		sqlitedb.Update(&node)
	}

	controllers.Response(c, common.OK, "", nil)
}

// 查看全部在线用户
func PageNodes(c *gin.Context) {

	node := models.Node{}
	web.GetModel(&node, c)

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

	total := sqlitedb.QueryPage(curPage, pageRecord, &nodes, "ip like ?", "%"+node.Ip+"%")

	data := make(map[string]interface{})

	data["nodes"] = nodes
	data["curPage"] = curPage
	data["pageRecord"] = pageRecord
	data["total"] = total
	controllers.Response(c, common.OK, "", data)
}
