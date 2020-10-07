/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 14:41
 */

package models

import "github.com/jinzhu/gorm"

type Page struct {
	CurPage  int `json:"curPage,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
}

//1对1聊天
type Node struct {
	gorm.Model
	Ip          string `json:"ip,omitempty"`
	Port        int    `json:"port,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	ClusterId   uint   `json:"clusterId,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
	NodeType    int    `json:"nodeType,omitempty"` //1 主机主节点，2 集群主节点，3，集群从节点
	Remark      string `json:"remark,omitempty"`
	Status      int    `json:"status,omitempty"` //1正常，2非正常
}

type Cluster struct {
	gorm.Model
	ClusterName string `json:"clusterName,omitempty"`
	Category    string `json:"category,omitempty"`
	Remark      string `json:"remark,omitempty"`
	Status      int    `json:"status,omitempty"`
}

type Project struct {
	gorm.Model
	ProjectName string `json:"clusterName,omitempty"`
	Remark      string `json:"remark,omitempty"`
}
