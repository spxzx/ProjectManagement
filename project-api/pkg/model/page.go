package model

import "github.com/gin-gonic/gin"

type Page struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"pageSize" form:"pageSize"`
}

func (p *Page) Bind(ctx *gin.Context) error {
	err := ctx.ShouldBind(&p)
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	if p.Page == 0 {
		p.Page = 1
	}
	return err
}
