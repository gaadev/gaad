package base

import (
	"encoding/json"
	"gaad/common"
	"gaad/db/sqlitedb"
	"gaad/models"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type QueryHandler func() (query interface{}, args []interface{})

type WhereHandler func() (where []interface{})

type CheckParam func(c *gin.Context) error

func Create(c *gin.Context, model interface{}, checkParam CheckParam) {
	createOrUpdate("create", c, model, checkParam)
}

func Delete(c *gin.Context, model interface{}) {
	err := GetModel(model, c)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}

	sqlitedb.Delete(model)
	controllers.Response(c, common.OK, "", nil)
}

func Update(c *gin.Context, model interface{}, checkParam CheckParam) {
	createOrUpdate("update", c, model, checkParam)
}

func Page(c *gin.Context, entity interface{}, entities interface{}, checkParam CheckParam, handler QueryHandler) {
	var (
		err error
	)

	page := models.Page{}

	dataByt, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(dataByt, entity)
	err = json.Unmarshal(dataByt, &page)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}

	if page.CurPage == 0 {
		page.CurPage = 1
	}
	if page.PageRecord == 0 {
		page.PageRecord = 10
	}
	if err = checkParam(c); err != nil {
		return
	}

	query, where := handler()
	total := sqlitedb.QueryPage(page.CurPage, page.PageRecord, entities, query, where...)

	data := make(map[string]interface{})

	data["data"] = entities
	data["curPage"] = page.CurPage
	data["pageRecord"] = page.PageRecord
	data["total"] = total
	controllers.Response(c, common.OK, "", data)
}

func createOrUpdate(opreation string, c *gin.Context, model interface{}, checkParam CheckParam) {
	err := GetModel(model, c)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}
	if err = checkParam(c); err != nil {
		return
	}
	if opreation == "create" {
		sqlitedb.Create(model)
	}
	if opreation == "update" {
		sqlitedb.Update(model)
	}

	controllers.Response(c, common.OK, "", nil)
}

func List(c *gin.Context, entities interface{}, handler WhereHandler) {
	where := handler()
	sqlitedb.QueryList(entities, where...)
	controllers.Response(c, common.OK, "", entities)
}
