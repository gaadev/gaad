package application

import (
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
// @Success 200 {object} common.JsonResult
// @Router /project/createProject [post]
// @Tags 项目(Project)
func CreateProject(c *gin.Context) {
	project := models.Project{}
	base.Create(c, &project, func(c *gin.Context) error {
		if project.ProjectName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		pro := models.Project{}
		sqlitedb.First(&pro, "ws_code = ?", pro.WsCode)
		//pro.Id > 0说明已经存在
		if pro.ID > 0 {
			return controllers.Response(c, common.OperationFailure, "WsCode重复", nil)
		}
		return nil
	})
}

// @Description 更新项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/updateProject [put]
// @Tags 项目(Project)
func UpdateProject(c *gin.Context) {
	project := models.Project{}
	base.Update(c, &project, func(c *gin.Context) error {
		if project.ProjectName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 删除项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/deleteProject [delete]
// @Tags 项目(Project)
func DeleteProject(c *gin.Context) {
	base.Delete(c, &models.Project{})
}

// @Description 分页查询项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/pageProjects [post]
// @Tags 项目(Project)
func PageProjects(c *gin.Context) {
	project := models.Project{}
	var projects []models.Project

	base.Page(c, &project, &projects,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 0)
			query = "project_name like ?"
			where = append(where, "%"+project.ProjectName+"%")
			return
		})
}

// @Description 查询所有项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/listProjects [post]
// @Tags 项目(Project)
func ListProjects(c *gin.Context) {
	var projects []models.Project
	base.List(c, &projects, func() (where []interface{}) {
		return nil
	})
}
