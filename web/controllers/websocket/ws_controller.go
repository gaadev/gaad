package websocket

import (
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
// @Success 200 {object} models.Rsp
// @Router /ws/shellConnect [get]
// @Tags Shell(shell)
func ShellConnect(c *gin.Context) {
	nodeId, err := strconv.ParseUint(c.Query("id"), 10, 64)
	cols, err := strconv.Atoi(c.Query("cols"))
	rows, err := strconv.Atoi(c.Query("rows"))
	if nil != err {
		controllers.Response(models.ParameterIllegal, "参数转换异常", nil).Write(c)
		return
	}
	node := models.Node{}
	//查询获取node
	sqlitedb.First(&node, " id = ?", nodeId)
	if nil == &node {
		controllers.Response(models.NotData, "连接不存在", nil).Write(c)
		return
	}
	base.HandleWsAndShell(&node, cols, rows, c)
	controllers.Response(models.OK, "", nil).Write(c)
}
