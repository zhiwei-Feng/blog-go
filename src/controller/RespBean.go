package controller

type RespBean struct {
	// http返回封装
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
