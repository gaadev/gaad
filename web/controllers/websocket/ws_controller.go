package websocket

import (
	"gaad/common"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"strconv"
)

// @Description shell连接
// @Param id query string true "节点id"
// @Param cols query int true "列数"
// @Param rows query int true "行数"
// @Success 200 {object} common.JsonResult
// @Router /ws/shellConnect [get]
// @Tags Shell(shell)
func ShellConnect(c *gin.Context) {
	nodeId, err := strconv.ParseUint(c.Query("id"), 10, 64)
	cols, err := strconv.Atoi(c.Query("cols"))
	rows, err := strconv.Atoi(c.Query("rows"))
	if nil != err {
		controllers.Response(c, common.ParameterIllegal, "参数转换异常", nil)
		return
	}
	node := models.Node{}
	//查询获取node
	sqlitedb.First(&node, " id = ?", nodeId)
	if nil == &node {
		controllers.Response(c, common.NotData, "连接不存在", nil)
		return
	}
	base.HandleWsAndShell(&node, cols, rows, c)
	return
}
