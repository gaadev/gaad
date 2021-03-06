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
// @Success 200 {object} models.Rsp
// @Router /service/createService [post]
// @Tags 服务(Service)
func CreateService(c *gin.Context) {
	service := models.Service{}
	rsp := base.Create(c, &service, func() *models.Rsp {
		if service.ServiceName == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		serv := models.Service{}
		sqlitedb.First(&serv, "service_code = ?", service.ServiceCode)
		//pro.Id > 0说明已经存在
		if serv.ID > 0 {
			return controllers.Response(models.OperationFailure, "服务标识重复", nil)
		}
		//系统内部查看project,填充资料
		proj := models.Project{}
		sqlitedb.First(&proj, "ID = ?", service.ProjectId)
		if proj.ID < 0 {
			return controllers.Response(models.OperationFailure, "项目不存在", nil)
		}
		service.ProjectName = proj.ProjectName
		service.WsCode = proj.WsCode
		if service.DevopsOpts == "" {
			service.DevopsOpts = "{}"
		}

		return nil
	})
	if rsp != nil {
		rsp.Write(c)
	} else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 更新项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} models.Rsp
// @Router /service/updateService [put]
// @Tags 服务(Service)
func UpdateService(c *gin.Context) {
	service := models.Service{}
	rsp := base.Update(c, &service, func() *models.Rsp {
		if service.ServiceName == "" {
			return controllers.Response(models.ParameterIllegal, "", nil)
		}
		if service.DevopsOpts == "" {
			service.DevopsOpts = "{}"
		}
		return nil
	})

	if rsp != nil {
		rsp.Write(c)
	} else {
		controllers.Response(models.OK, "", nil).Write(c)
	}

}

// @Description 删除项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} models.Rsp
// @Router /service/deleteService [delete]
// @Tags 服务(Service)
func DeleteService(c *gin.Context) {
	rsp := base.Delete(c, &models.Service{})
	if rsp != nil {
		rsp.Write(c)
	} else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

// @Description 分页查询项目
// @Accept  json
// @Produce json
// @Param data body models.Service true "Data"
// @Success 200 {object} models.Rsp
// @Router /service/pageServices [post]
// @Tags 服务(Service)
func PageServices(c *gin.Context) {
	service := models.Service{}
	var services []models.Service

	rsp := base.Page(c, &service, &services,
		func() *models.Rsp {
			return nil
		},
		func() (query interface{}, where []interface{}) {
			where = make([]interface{}, 0)
			sql := "1 = 1 "
			if service.ServiceName != "" {
				sql += " and service_name like ? "
				where = append(where, "%"+service.ServiceName+"%")
			}
			if service.ProjectId != 0 {
				sql += " and project_id = ? "
				where = append(where, service.ProjectId)
			}
			query = sql

			return
		})
	if rsp != nil {
		rsp.Write(c)
	} else {
		controllers.Response(models.OK, "", nil).Write(c)
	}
}

func toDevopsOpt(buildOpt string) string {
	switch buildOpt {
	case "JavaOpts":
		return "--java-opts"
	case "BuildTool":
		return "--build-tool"
	case "Workspace":
		return "--workspace"
	case "Dockerfile":
		return "--dockerfile"
	case "Template":
		return "--template"
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
	conCatOptFeildAndValue(&build, toDevopsOpt("GitUrl"), service.GitUrl)
	conCatOptFeildAndValue(&build, toDevopsOpt("GitBranch"), service.GitBranch)
	build.WriteString(GenDevopsOpts(service))
	build.WriteString(" ")
	build.WriteString(service.ServiceCode)
	cmd = build.String()
	return
}

func GenDevopsOpts(service *models.Service) string {
	devopsOpts := models.DevopsOpts{}
	json.Unmarshal([]byte(service.DevopsOpts), &devopsOpts)
	var build strings.Builder

	t := reflect.TypeOf(devopsOpts)
	v := reflect.ValueOf(devopsOpts)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).String() != "" {
			conCatOptFeildAndValue(&build, toDevopsOpt(t.Field(k).Name), v.Field(k).String())
		}
	}
	if service.Dockerfile != "" {
		conCatOptFeildAndValue(&build, toDevopsOpt("Dockerfile"), service.Dockerfile)
	}
	if service.Template != "" {
		conCatOptFeildAndValue(&build, toDevopsOpt("Template"), service.Template)
	}
	conCatOptFeildAndValue(&build, toDevopsOpt("Workspace"), service.WsCode)

	OptStr := build.String()
	return OptStr

}

func conCatOptFeildAndValue(build *strings.Builder, key string, value string) {
	build.WriteString(" ")
	build.WriteString(key)
	build.WriteString(" '")
	build.WriteString(value)
	build.WriteString("' ")
}

// @Description 部署的接口还未完善
// @Accept  json
// @Produce json
// @Success 200 {object} models.Rsp
// @Router /service/deploy [post]
// @Tags 服务(Service)
func Deploy(c *gin.Context) {

	var (
		devopsCmd string
		status    int
	)
	service := &models.Service{}

	err := base.GetModel(&service, c)
	if err != nil {
		controllers.Response(models.ParameterIllegal, "参数格式有误", nil).Write(c)
		return
	}

	serv := &models.Service{}
	sqlitedb.First(serv, "id = ? and status = 1", service.ID)
	if serv.ID != 0 {

		devopsCmd = GenDevopsCmd(serv)

	} else {
		controllers.Response(models.OperationFailure, "操作失败，服务未启用", nil).Write(c)
		return
	}

	go func() {

		par := []string{
			"-c",
			"cd ./script/devops/bin; ./" + devopsCmd,
		}

		num:= serv.DeployNum
		boltdb.Update(service.WsCode+":"+service.ServiceCode, strconv.Itoa(num+1))

		logPath := GetLogPath(service)
		common.CreateFile(logPath)
		logFilePath := logPath + strconv.Itoa(num+1) + ".log"

		serv.DeployNum = num + 1
		sqlitedb.Update(serv)

		retErr := common.DeployCommand(service, logFilePath, "/bin/sh", par)

		if retErr == nil {
			status = 1
		} else {
			status = 2
		}

		deploy := &models.Deploy{ServiceId: service.ID, ServiceName: service.ServiceName,
			ServiceCode: service.ServiceCode, DeployNum: num + 1, LogFilePath: logFilePath, Status: status}
		sqlitedb.Create(deploy)

		data := make(map[string]interface{})

		data["status"] = "ok"

	}()

	controllers.Response(models.OK, "开始构建", "").Write(c)
}

// @Description 查询当前服务的deploy
// @Accept  json
// @Produce json
// @Param  serviceId query  string  true "serviceId"
// @Success 200 {object} models.Rsp
// @Router /service/listDevops [get]
// @Tags 服务(Service)
func ListDevops(c *gin.Context) {
	serviceId := c.Query("serviceId")
	if serviceId == "" {
		controllers.Response(models.OperationFailure, "serviceId不能为空", nil).Write(c)
	}
	var deploys []models.Deploy
	rsp := base.List( &deploys, func() (where []interface{}) {

		where = make([]interface{}, 0)

		sql := "service_id = ? "
		where = append(where, sql)
		where = append(where,serviceId)
		return
	})
	if rsp != nil && rsp.Code != 200 {
		rsp.Write(c)
		return
	}
	rsp.Write(c)
}

func GetLogPath(service *models.Service) string {
	return "./script/devops/workspace/" + service.WsCode + "/log/" + service.ServiceCode + "/"
}

func getServiceDevopsNumKey(service *models.Service) string {
	return service.WsCode + ":" + service.ServiceCode + ":" + "devopsNum"
}



// @Description 展示
// @Accept  json
// @Produce json
// @Param  serviceId query  string  true "serviceId"
// @Param  logNum query  string  false "logNum"
// @Param  start query  string  false "start"
// @Router /service/display [get]
// @Tags 服务(Service)
func Display(c *gin.Context) {
	var (
		count   int
		builder strings.Builder
	)
	serviceId := c.Query("serviceId")
	logNum := c.Query("logNum")
	startStr := c.Query("start")

	serv :=  &models.Service{}
	sqlitedb.GetById(serv,serviceId)
	if startStr == "" {
		startStr = "1"
	}
	start, err := strconv.Atoi(startStr)

	flag.Parse()

	logPath := GetLogPath(serv)
	if logNum == "" {
		logNum = "1"
	}
	if serv.DeployNum != 0  {
		logNum =  strconv.Itoa(serv.DeployNum)
	}
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
	count++
	c.Writer.Header().Add("Content-Type","text/html;charset=utf-8")
	c.Writer.Header().Add("X-Text-Lines", strconv.Itoa(count))
	c.Writer.Write([]byte(builder.String()))

}
