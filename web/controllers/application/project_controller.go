package application

import (
	"gaad/common"
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
func CreateProject(c *gin.Context) {
	project := models.Project{}
	base.Create(c, &project, func(c *gin.Context) error {
		if project.ProjectName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
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
func DeleteProject(c *gin.Context) {
	base.Delete(c, &models.Project{})
}

// @Description 分页查询项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/pageProjects [post]
func PageProjects(c *gin.Context) {
	project := models.Project{}
	var projects []models.Project

	base.Page(c, &project, &projects,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 3)
			query = "project_name like ?"
			where[0] = "%" + project.ProjectName + "%"
			return
		})
}

// @Description 查询所有项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /project/listProjects [post]
func ListProjects(c *gin.Context) {
	var projects []models.Project
	base.List(c, &projects, func() (where []interface{}) {
		return nil
	})
}