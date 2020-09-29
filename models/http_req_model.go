/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 14:41
 */

package models

import "github.com/jinzhu/gorm"

//1对1聊天
type Node struct {
	gorm.Model
	Ip        string `json:"ip,omitempty"`
	Port      int    `json:"port,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	ClusterId string `json:"clusterId,omitempty"`
	Remark    string `json:"remark,omitempty"`
}
