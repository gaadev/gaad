package node

import (
	"fmt"
	"gaad/common"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"net"
)

// 查看全部在线用户
func CreateNode(c *gin.Context) {
	node := models.Node{}
	base.Create(c, &node, func(c *gin.Context) error {
		address := net.ParseIP(node.Ip)
		if address == nil {
			return controllers.Response(c, common.ParameterIllegal, "Ip地址格式有误", nil)
		} else {
			fmt.Println("正确的ip地址", address.String())
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

// 查看全部在线用户
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

// 查看全部在线用户
func DeleteNode(c *gin.Context) {

	base.Delete(c, &models.Node{})
}

// 查看全部在线用户
func PageNodes(c *gin.Context) {
	node := models.Node{}
	var nodes []models.Node

	base.Page(c, &node, &nodes, func() (query interface{}, where []interface{}) {
		where = make([]interface{}, 3)
		query = "ip like ?"
		where[0] = "%" + node.Ip + "%"
		return
	})
}
