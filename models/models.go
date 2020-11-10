/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 14:41
 */

package models

import "github.com/jinzhu/gorm"

type Page struct {
	CurPage    int `json:"curPage,omitempty"`
	PageRecord int `json:"pageRecord,omitempty"`
}

//节点
type Node struct {
	gorm.Model 		   `swaggerignore:"true"`
	Ip          string `json:"ip,omitempty"`          //ip地址
	Port        int    `json:"port,omitempty"`        //端口号
	Username    string `json:"username,omitempty"`    //用户名
	Password    string `json:"password,omitempty"`    //密码
	ClusterId   uint   `json:"clusterId,omitempty"`   //集群id
	ClusterName string `json:"clusterName,omitempty"` //集群名称
	NodeType    int    `json:"nodeType,omitempty"`    //1 主机主节点，2 集群主节点，3，集群从节点
	Remark      string `json:"remark,omitempty"`      //标记
	Status      int    `json:"status,omitempty"`      //状态：1正常，2非正常
}

//集群
type Cluster struct {
	gorm.Model         `swaggerignore:"true"`
	ClusterName string `json:"clusterName,omitempty"`
	Category    string `json:"category,omitempty"`
	Remark      string `json:"remark,omitempty"`
	Status      int    `json:"status,omitempty"`
}

//项目
type Project struct {
	gorm.Model         `swaggerignore:"true"`
	ProjectName string `json:"projectName,omitempty"`
	WsCode      string `json:"wsCode,omitempty"`
	ClusterId   uint   `json:"clusterId,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
	Status      int    `json:"status,omitempty"`
	Remark      string `json:"remark,omitempty"`
	GitAccount  string `json:"gitAccount,omitempty"`
	GitPassword string `json:"gitPassword,omitempty"`
}

//服务
type Service struct {
	gorm.Model         `swaggerignore:"true"`
	ServiceName string `json:"serviceName,omitempty"`
	ServiceCode string `json:"serviceCode,omitempty"`
	GitUrl      string `json:"gitUrl,omitempty"`
	GitBranch   string `json:"gitBranch,omitempty"`
	Lang        string `json:"lang,omitempty"`
	DevopsOpts	string `json:"devopsOpts,omitempty"`
	ProjectId   uint   `json:"projectId,omitempty"`
	WsCode      string `json:"wsCode,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
	Dockerfile  string `json:"dockerfile,omitempty"`
	Template    string `json:"template,omitempty"`
	DeployNum   int `json:"deployNum,omitempty"`
	Status      int    `json:"status,omitempty"`
	Remark      string `json:"remark,omitempty"`
}

type Deploy struct {
	gorm.Model         `swaggerignore:"true"`
	ServiceId   uint   `json:"serviceId,omitempty"`
	ServiceName string `json:"serviceName,omitempty"`
	ServiceCode string `json:"serviceCode,omitempty"`
	DeployNum   int    `json:"deployNum,omitempty"`
	Status      int    `json:"status,omitempty"` //1成功，2失败
	LogFilePath string `json:"logFilePath,omitempty"`
}

type DevopsOpts struct {
	JavaOpts  string `json:"javaOpts,omitempty"`
	BuildTool string `json:"buildTool,omitempty"`
}
