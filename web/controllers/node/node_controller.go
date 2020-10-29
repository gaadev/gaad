package node

import (
	"fmt"
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
// @Success 200 {object} models.Rsp
// @Router /node/createNode [post]
// @Tags 节点(node)
func CreateNode(c *gin.Context) {
	node := models.Node{}


	rsp := base.Create(c, &node, func() *models.Rsp {
		address := net.ParseIP(node.Ip)
		if address == nil {
			return controllers.Response(models.ParameterIllegal, "Ip地址格式有误", nil)
		} else {
			fmt.Println("正确的ip地址", address.String())
		}

		nod := models.Node{}
		sqlitedb.First(&nod, "ip = ?", nod.Ip)
		//pro.Id > 0说明已经存在
		if nod.ID > 0 {
			return controllers.Response(models.OperationFailure, "IP重复", nil)
		}

		if node.Port < 2 {
			return controllers.Response(models.ParameterIllegal, "端口不可以小于2", nil)
		}

		if node.Username == "" || node.Password == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		return nil
	})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 更新主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /node/updateNode [put]
// @Tags 节点(node)
func UpdateNode(c *gin.Context) {
	node := models.Node{}
	rsp := base.Update(c, &node, func() *models.Rsp {
		if node.Ip != "" {
			address := net.ParseIP(node.Ip)
			if address == nil {
				return controllers.Response(models.ParameterIllegal, "Ip地址格式有误", nil)
			} else {
				fmt.Println("正确的ip地址", address.String())
			}
		}

		if node.Port != 0 && node.Port < 2 {
			return controllers.Response(models.ParameterIllegal, "端口不可以小于2", nil)
		}
		return nil
	})


	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}

}

// @Description 删除主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /node/deleteNode [delete]
// @Tags 节点(node)
func DeleteNode(c *gin.Context) {
	rsp := base.Delete(c, &models.Node{})
	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 查寻主机节点
// @Accept  json
// @Produce json
// @Param data body models.Node true "Data"
// @Success 200 {object} models.Rsp
// @Router /node/pageNodes [post]
// @Tags 节点(node)
func PageNodes(c *gin.Context) {

	node := models.Node{}
	var nodes []models.Node

	rsp := base.Page(c, &node, &nodes,
		func() *models.Rsp {
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
	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}
