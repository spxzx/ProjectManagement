package ath

type ProjectAuthReq struct {
	Action string `form:"action"`
	Id     int64  `form:"id"`
	Nodes  string `form:"nodes" json:"nodes"`
}
