package web

import (
	"encoding/json"
	"gaad/common"
	"gaad/controllers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func GetModel(model interface{}, c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, model)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "参数格式有误", nil)
		return
	}
}
