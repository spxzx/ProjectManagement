package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model/member"
	"github.com/spxzx/project-api/pkg/model/poj"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/account/account"
	"net/http"
	"time"
)

type HandlerAccount struct{}

func NewAccount() *HandlerAccount {
	return &HandlerAccount{}
}

func (*HandlerAccount) getAccountList(ctx *gin.Context) {
	var req *member.AccountReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	actResp, err := rpc.AccountServiceClient.GetAccountList(c, &account.AccountRpcRequest{
		MemberId:         ctx.GetInt64("memberId"),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		OrganizationCode: ctx.GetString("organizationCode"),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*member.Account
	_ = copier.Copy(&list, actResp.AccountList)
	if list == nil || len(list) <= 0 {
		list = []*member.Account{}
	}
	var authList []*poj.ProjectAuth
	_ = copier.Copy(&authList, actResp.AuthList)
	if authList == nil {
		authList = []*poj.ProjectAuth{}
	}
	r.Success(ctx, gin.H{
		"total":    actResp.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	})
}
