package project

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model"
	"github.com/spxzx/project-api/pkg/model/ath"
	"github.com/spxzx/project-api/pkg/model/poj"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/auth/auth"
	"net/http"
	"time"
)

type HandlerAuth struct{}

func NewAuth() *HandlerAuth {
	return &HandlerAuth{}
}

func (*HandlerAuth) getAuthList(ctx *gin.Context) {
	var page = &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.AuthServiceClient.AuthList(c, &auth.AuthRpcRequest{
		OrganizationCode: ctx.GetString("organizationCode"),
		Page:             page.Page,
		PageSize:         page.PageSize,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var authList []*poj.ProjectAuth
	_ = copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*poj.ProjectAuth{}
	}
	r.Success(ctx, gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	})
}

func (*HandlerAuth) apply(ctx *gin.Context) {
	var req *ath.ProjectAuthReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	var nodes []string
	if req.Nodes != "" {
		_ = json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	applyResponse, err := rpc.AuthServiceClient.Apply(c, &auth.AuthRpcRequest{
		Action: req.Action,
		Id:     req.Id,
		Nodes:  nodes,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*poj.ProjectNodeAuthTree
	_ = copier.Copy(&list, applyResponse.List)
	var checkedList []string
	_ = copier.Copy(&checkedList, applyResponse.CheckedList)
	r.Success(ctx, gin.H{
		"list":        list,
		"checkedList": checkedList,
	})
}

func (*HandlerAuth) GetAuthNodes(ctx *gin.Context) ([]string, error) {
	memberId := ctx.GetInt64("memberId")
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.AuthServiceClient.AuthNodesByMemberId(c, &auth.AuthRpcRequest{
		MemberId: memberId,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return nil, err
	}
	return response.List, err
}
