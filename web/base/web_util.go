package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func GetModel(model interface{}, c *gin.Context) error {
	data, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, model)
	return err
}
