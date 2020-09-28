/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 09:59
 */

package main

import (
	"gaad/initialize"
	"gaad/routers"
)

func main() {
	initialize.InitConfig()

	initialize.InitFile()
	//初始化http路由，并启动http服务器
	routers.InitWebRouters()

	//启动http服务器，阻塞线程
	routers.InitHttpServer()
}