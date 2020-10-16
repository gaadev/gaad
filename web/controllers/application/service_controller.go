package application

import (
	"bufio"
	"encoding/json"
	"flag"
	"gaad/common"
	"gaad/db/boltdb"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/base"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// @Description 创建项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} common.JsonResult
// @Router /service/createService [post]
// @Tags 服务(Service)
func CreateService(c *gin.Context) {
	service := models.Service{}
	base.Create(c, &service, func(c *gin.Context) error {
		if service.ServiceName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		serv := models.Service{}
		sqlitedb.First(&serv, "service_code = ?", service.ServiceCode)
		//pro.Id > 0说明已经存在
		if serv.ID > 0 {
			return controllers.Response(c, common.OperationFailure, "服务标识重复", nil)
		}
		//系统内部查看project,填充资料
		proj := models.Project{}
		sqlitedb.First(&proj, "ID = ?", serv.ProjectId)
		if proj.ID < 0 {
			return controllers.Response(c, common.OperationFailure, "项目不存在", nil)
		}
		serv.ProjectName = proj.ProjectName
		serv.WsCode = proj.WsCode
		return nil
	})
}

// @Description 更新项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} common.JsonResult
// @Router /service/updateService [put]
// @Tags 服务(Service)
func UpdateService(c *gin.Context) {
	service := models.Service{}
	base.Update(c, &service, func(c *gin.Context) error {
		if service.ServiceName == "" {
			return controllers.Response(c, common.ParameterIllegal, "", nil)
		}
		return nil
	})

}

// @Description 删除项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} common.JsonResult
// @Router /service/deleteService [delete]
// @Tags 服务(Service)
func DeleteService(c *gin.Context) {
	base.Delete(c, &models.Service{})
}

// @Description 分页查询项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} common.JsonResult
// @Router /service/pageServices [post]
// @Tags 服务(Service)
func PageServices(c *gin.Context) {
	service := models.Service{}
	var services []models.Service

	base.Page(c, &service, &services,
		func(c *gin.Context) error {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 0)
			sql := "1 = 1"
			if service.ServiceName != "" {
				query = "service_name like ?"
				where = append(where, "%"+service.ServiceName+"%")
			}
			query = sql

			return
		})
}

func toDevopsOpt(buildOpt string) string {
	switch buildOpt {
	case "JavaOpts":
		return "--java-opts"
	case "BuildTool":
		return "--build-tool"
	case "Workspace":
		return "--workspace"
	case "GitUrl":
		return "--git-url"
	case "GitBranch":
		return "--git-branch"
	default:
		return "--help"
	}
}

func GenDevopsCmd(service *models.Service) (cmd string) {
	var build strings.Builder

	build.WriteString("devops run ")
	build.WriteString(service.Lang)
	build.WriteString(" ")
	conCatOptFeildAndValue(build, toDevopsOpt("GitUrl"), service.GitUrl)
	conCatOptFeildAndValue(build, toDevopsOpt("GitBranch"), service.GitBranch)
	build.WriteString(GenDevopsOpts(service))
	build.WriteString(" ")
	build.WriteString(service.ServiceCode)
	cmd = build.String()
	return
}

func GenDevopsOpts(service *models.Service) string {

	devopsOpts := models.DevopsOpts{}
	err := json.Unmarshal([]byte(service.DevopsOpts), &devopsOpts)
	if err != nil {
		log.Fatal(err)
	}
	var build strings.Builder

	t := reflect.TypeOf(devopsOpts)
	v := reflect.ValueOf(devopsOpts)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).String() != "" {
			conCatOptFeildAndValue(build, toDevopsOpt(t.Field(k).Name), v.Field(k).String())
		}
	}
	conCatOptFeildAndValue(build, toDevopsOpt("Workspace"), service.WsCode)

	OptStr := build.String()
	return OptStr

}

func conCatOptFeildAndValue(build strings.Builder, key string, value string) {
	build.WriteString(" ")
	build.WriteString(key)
	build.WriteString("='")
	build.WriteString(value)
	build.WriteString("' ")
}

// @Description 部署的接口还未完善
// @Accept  json
// @Produce json
// @Success 200 {object} common.JsonResult
// @Router /service/deploy [post]
// @Tags 服务(Service)
func Deploy(c *gin.Context) {

	var (
		devopsCmd string
		status    int
	)
	service := models.Service{}

	err := base.GetModel(&service, c)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}

	serv := &models.Service{}
	sqlitedb.First(serv, "id = ? and status = 1", service.ID)
	if serv.ID != 0 {

		devopsCmd = GenDevopsCmd(serv)

	} else {
		controllers.Response(c, common.OperationFailure, "操作失败，服务未启用", nil)
		return
	}

	par := []string{
		"-c",
		"cd ./script/devops/bin; ./" + devopsCmd,
	}

	devopsNum := boltdb.View(getServiceDevopsNumKey(service))
	num, err := strconv.Atoi(devopsNum)
	boltdb.Update(service.WsCode+":"+service.ServiceCode, strconv.Itoa(num+1))

	logPath := GetLogPath(service)
	logFilePath := logPath + strconv.Itoa(num+1) + ".log"
	common.CreateFile(logFilePath)

	retErr := common.DeployCommand(service, logFilePath, "/bin/sh", par)

	if retErr == nil {
		status = 1
	} else {
		status = 2
	}

	deploy := &models.Deploy{ServiceId: service.ID, ServiceName: service.ServiceName,
		ServiceCode: service.ServiceCode, DeployNum: strconv.Itoa(num + 1), LogFilePath: logFilePath, Status: status}
	sqlitedb.Create(deploy)

	data := make(map[string]interface{})

	data["status"] = "ok"

	controllers.Response(c, common.OK, "", data)

}

// @Description 查询所有项目
// @Accept  json
// @Produce json
// @Param data body models.Project true "Data"
// @Success 200 {object} common.JsonResult
// @Router /service/ListDevops [post]
// @Tags 服务(Service)
func ListDevops(c *gin.Context) {
	var deploys []models.Deploy
	base.List(c, &deploys, func() (where []interface{}) {
		return nil
	})

}

func GetLogPath(service models.Service) string {
	return "./script/devops/workspace/" + service.WsCode + "/log/" + service.ServiceCode + "/"
}

func getServiceDevopsNumKey(service models.Service) string {
	return service.WsCode + ":" + service.ServiceCode + ":" + "devopsNum"
}

// @Description 展示
// @Accept  json
// @Produce json
// @Success 200 {object} common.JsonResult
// @Router /service/display [post]
// @Tags 服务(Service)
func Display(c *gin.Context) {
	var (
		count   int
		builder strings.Builder
	)
	service := models.Service{}

	err := base.GetModel(&service, c)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}
	logNum := c.Query("logNum")
	startStr := c.Query("start")
	if startStr == "" {
		startStr = "1"
	}
	start, err := strconv.Atoi(startStr)

	flag.Parse()

	logPath := GetLogPath(service)
	f, err := os.Open(logPath + logNum + ".log")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)
	count += start
	for i := 0; i < start; i++ {
		s.Scan()
	}

	for s.Scan() {
		count++
		builder.WriteString(s.Text())
		builder.WriteString("\n")
	}
	//strings.LastIndex(deployLog,"")
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]interface{})

	data["data"] = builder.String()
	data["end"] = count

	controllers.Response(c, common.OK, "", data)

}
