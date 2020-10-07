package application

import (
	"bufio"
	"flag"
	"gaad/common"
	"gaad/web/controllers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"strings"
)

// @Description 部署的接口还未完善
// @Accept  json
// @Produce json
// @Success 200 {object} common.JsonResult
// @Router /service/deploy [post]
// @Tags 部署(Deploy)
func Deploy(c *gin.Context) {
	//for i := 0; i < 100; i++{
	//	time.Sleep(time.Duration(1)*time.Second)
	//	f.WriteString("Hello world\n")
	//}
	par := []string{
		"-c",
		"devops run java --git-url 'http://192.168.10.235/unsun/biz/jshyun-console.git' --git-branch test --java-opts '-Dprofile=test -Dconfig-registry=core-config' --workspace unsun  console-advertisement",
	}
	common.ExecCommand("/bin/sh", par)

	data := make(map[string]interface{})

	data["status"] = "ok"

	controllers.Response(c, common.OK, "", data)

}

// @Description 展示
// @Accept  json
// @Produce json
// @Success 200 {object} common.JsonResult
// @Router /service/display [post]
// @Tags 部署(Deploy)
func Display(c *gin.Context) {
	var (
		count   int
		builder strings.Builder
	)
	startStr := c.Query("start")
	if startStr == "" {
		startStr = "1"
	}
	start, err := strconv.Atoi(startStr)

	flag.Parse()

	f, err := os.Open("./log/log-1.log")
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
