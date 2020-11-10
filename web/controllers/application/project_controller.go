package application

import (
	"encoding/base64"
	"fmt"
	"gaad/common"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
)

// @Description 创建项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} models.Rsp
// @Router /project/createProject [post]
// @Tags 项目(Project)
func CreateProject(c *gin.Context) {
	project := models.Project{}
	rsp := base.Create(c, &project, func() *models.Rsp {
		if project.ProjectName == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		pro := models.Project{}
		sqlitedb.First(&pro, "ws_code = ?", project.WsCode)
		//pro.Id > 0说明已经存在
		if pro.ID > 0 {
			return controllers.Response(models.OperationFailure, "WsCode重复", nil)
		}
		return nil
	})
	if rsp != nil && rsp.Code != 200 {
		rsp.Write(c)
		return
	}

	cluster := &models.Cluster{}
	sqlitedb.First(cluster, "ID = ? and status = 1", project.ClusterId)

	if cluster.ID == 0 {
		controllers.Response(models.OperationFailure, "关联集群状态不正常", nil).Write(c)
		return
	}

	var (
		nodes       []models.Node
		nodeMasters []models.Node
	)
	sqlitedb.QueryList(&nodes, "cluster_id = ?", cluster.ID)

	for _, node := range nodes {
		if node.NodeType == 2 {
			nodeMasters = append(nodeMasters, node)
		}
	}

	if len(nodeMasters) == 0 {
		controllers.Response(models.OperationFailure, "关联集群没有主节点", nil).Write(c)
		return
	}

	nodeMaster := nodeMasters[0]

	nodeSecretMsg := nodeMaster.Username + ":" + nodeMaster.Ip + ":" + nodeMaster.Password
	secret := base64.StdEncoding.EncodeToString([]byte(nodeSecretMsg))

	remoteDeploySecret := project.WsCode + "=" + secret
	fmt.Println(remoteDeploySecret)

	par := []string{
		"-c",
		"sh ./shell/add_node_secret.sh " + project.WsCode + " " + remoteDeploySecret,
	}

	common.ExecCommand("/bin/sh", par)

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 更新项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} models.Rsp
// @Router /project/updateProject [put]
// @Tags 项目(Project)
func UpdateProject(c *gin.Context) {
	project := models.Project{}
	rsp := base.Update(c, &project, func() *models.Rsp {
		if project.ProjectName == "" {
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

// @Description 删除项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} models.Rsp
// @Router /project/deleteProject [delete]
// @Tags 项目(Project)
func DeleteProject(c *gin.Context) {
	rsp := base.Delete(c, &models.Project{})
	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 分页查询项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} models.Rsp
// @Router /project/pageProjects [post]
// @Tags 项目(Project)
func PageProjects(c *gin.Context) {
	project := models.Project{}
	var projects []models.Project

	rsp := base.Page(c, &project, &projects,
		func() *models.Rsp {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 0)
			sql := "1 = 1"
			if project.ProjectName != "" {
				query = "project_name like ?"
				where = append(where, "%"+project.ProjectName+"%")
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

// @Description 查询所有项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} models.Rsp
// @Router /project/listProjects [post]
// @Tags 项目(Project)
func ListProjects(c *gin.Context) {
	var projects []models.Project
	rsp := base.List( &projects, func() (where []interface{}) {
		where = make([]interface{}, 0)
		where = append(where, "status = 1")
		return
	})

	if rsp != nil {
		rsp.Write(c)
	}  else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}
