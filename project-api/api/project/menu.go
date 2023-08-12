package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	menu_ "github.com/spxzx/project-api/pkg/model/menu"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type HandlerMenu struct{}

func NewMenu() *HandlerMenu {
	return &HandlerMenu{}
}

func (*HandlerMenu) getMenuList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res, err := rpc.MenuServiceClient.GetMenuList(c, &emptypb.Empty{})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*menu_.Menu
	_ = copier.Copy(&list, res.List)
	if list == nil {
		list = []*menu_.Menu{}
	}
	r.Success(ctx, list)
}
