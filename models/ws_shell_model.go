package models

//websocket消息
type WsMsg struct {
	Type string `json:"type,omitempty"` //类型
	Cmd  string `json:"cmd,omitempty"`  //执行命令
	Cols int    `json:"cols,omitempty"` //列数
	Rows int    `json:"rows,omitempty"` //行数
}
