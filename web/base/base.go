package base

import (
	"encoding/json"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type QueryHandler func() (query interface{}, args []interface{})

type WhereHandler func() (where []interface{})

type CheckParam func() *models.Rsp

func Create(c *gin.Context, model interface{}, checkParam CheckParam) *models.Rsp {
	return createOrUpdate("create", c, model, checkParam)
}

func Delete(c *gin.Context, model interface{}) *models.Rsp  {
	err := GetModel(model, c)
	if err != nil {
		return controllers.Response(models.ParameterIllegal, "参数格式有误", nil)
	}

	sqlitedb.Delete(model)
	return controllers.Response(models.OK, "", nil)
}

func Update(c *gin.Context, model interface{}, checkParam CheckParam) *models.Rsp {
	return createOrUpdate("update", c, model, checkParam)
}

func Page(c *gin.Context, entity interface{}, entities interface{}, checkParam CheckParam, handler QueryHandler) *models.Rsp {
	var (
		err error
	)

	page := models.Page{}

	dataByt, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(dataByt, entity)
	err = json.Unmarshal(dataByt, &page)
	if err != nil {
		return controllers.Response(models.ParameterIllegal, "参数格式有误", nil)

	}

	if page.CurPage == 0 {
		page.CurPage = 1
	}
	if page.PageRecord == 0 {
		page.PageRecord = 10
	}
	if rsp := checkParam(); rsp != nil {
		return rsp
	}

	query, where := handler()
	total := sqlitedb.QueryPage(page.CurPage, page.PageRecord, entities, query, where...)

	data := make(map[string]interface{})

	data["data"] = entities
	data["curPage"] = page.CurPage
	data["pageRecord"] = page.PageRecord
	data["total"] = total

	return controllers.Response(models.OK, "", data)
}

func createOrUpdate(opreation string, c *gin.Context, model interface{}, checkParam CheckParam) *models.Rsp {
	err := GetModel(model, c)
	if err != nil {
		return controllers.Response(models.ParameterIllegal, "参数格式有误", nil)

	}
	if rsp := checkParam(); rsp != nil {
		return rsp
	}
	if opreation == "create" {
		sqlitedb.Create(model)
	}
	if opreation == "update" {
		sqlitedb.Update(model)
	}
	return controllers.Response(models.OK, "", nil)
}

func List(c *gin.Context, entities interface{}, handler WhereHandler) *models.Rsp  {
	where := handler()
	sqlitedb.QueryList(entities, where...)
	return controllers.Response(models.OK, "", entities)
}
