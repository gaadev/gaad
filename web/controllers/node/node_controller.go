package node

import (
	"fmt"
	"gaad/common"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"net"
)

// @Description 创建主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/createNode [post]
// @Tags 节点(node)
func CreateNode(c *gin.Context) {
	node := models.Node{}
	base.Create(c, &node, func(c *gin.Context) error {
		address := net.ParseIP(node.Ip)
		if address == nil {
			return controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
		} else {
			fmt.Println("正确的ip地址", address.String())
		}

		nod := models.Node{}
		sqlitedb.First(&nod, "ip = ?", nod.Ip)
		//pro.Id > 0说明已经存在
		if nod.ID > 0 {
			return controllers.Response(c, common.OperationFailure, "IP重复", nil)
		}

		if node.Port < 2 {
			return controllers.Response(c, common.ParameterIllegal, "端口不可以小于2", nil)
		}

		if node.Username == "" || node.Password == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 更新主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/updateNode [put]
// @Tags 节点(node)
func UpdateNode(c *gin.Context) {
	node := models.Node{}
	base.Update(c, &node, func(c *gin.Context) error {
		if node.Ip != "" {
			address := net.ParseIP(node.Ip)
			if address == nil {
				return controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
			} else {
				fmt.Println("正确的ip地址", address.String())
			}
		}

		if node.Port != 0 && node.Port < 2 {
			return controllers.Response(c, common.ParameterIllegal, "端口不可以小于2", nil)
		}
		return nil
	})

}

// @Description 删除主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/deleteNode [delete]
// @Tags 节点(node)
func DeleteNode(c *gin.Context) {

	base.Delete(c, &models.Node{})
}

// @Description 查寻主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} common.JsonResult
// @Router /node/pageNodes [post]
// @Tags 节点(node)
func PageNodes(c *gin.Context) {

	node := models.Node{}
	var nodes []models.Node

	base.Page(c, &node, &nodes,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, args []interface{}) {
			args = make([]interface{}, 0)

			sql := "1 = 1"
			if node.Ip != "" {
				sql += " and ip like ?"
				args = append(args, "%"+node.Ip+"%")
			}
			query = sql
			return
		})
}
