package application

import (
	"gaad/common"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
)

// 查看全部在线用户
func CreateProject(c *gin.Context) {
	project := models.Project{}
	base.Create(c, &project, func(c *gin.Context) error {
		if project.ProjectName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})
}

// 查看全部在线用户
func UpdateProject(c *gin.Context) {
	project := models.Project{}
	base.Update(c, &project, func(c *gin.Context) error {
		if project.ProjectName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// 查看全部在线用户
func DeleteProject(c *gin.Context) {
	base.Delete(c, &models.Project{})
}

// 查看全部在线用户
func PageProjects(c *gin.Context) {
	project := models.Project{}
	var projects []models.Project

	base.Page(c, &project, &projects, func() (query interface{}, where []interface{}) {
		where = make([]interface{}, 3)
		query = "project_name like ?"
		where[0] = "%" + project.ProjectName + "%"
		return
	})
}

func ListProjects(c *gin.Context) {
	var projects []models.Project
	base.List(c, &projects, func() (query interface{}, args []interface{}) {
		return nil, nil
	})
}
