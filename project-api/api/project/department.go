package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model/dpt"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/department/department"
	"github.com/spxzx/project-project/pkg/model"
	"net/http"
	"time"
)

type HandlerDepartment struct{}

func NewDepartment() *HandlerDepartment {
	return &HandlerDepartment{}
}

func (*HandlerDepartment) getDepartmentList(ctx *gin.Context) {
	var req *dpt.DepartmentReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dptResp, err := rpc.DepartmentServiceClient.GetDepartmentList(c,
		&department.DepartmentRpcRequest{
			Page:                 req.Page,
			PageSize:             req.PageSize,
			ParentDepartmentCode: req.Pcode,
			OrganizationCode:     ctx.GetString("organizationCode"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*dpt.Department
	_ = copier.Copy(&list, dptResp.List)
	if list == nil || len(list) <= 0 {
		list = []*dpt.Department{}
	}
	r.Success(ctx, gin.H{
		"total": dptResp.Total,
		"page":  req.Page,
		"list":  list,
	})
}

func (*HandlerDepartment) saveDepartment(ctx *gin.Context) {
	var req *dpt.DepartmentReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	departmentMessage, err := rpc.DepartmentServiceClient.SaveDepartment(c,
		&department.DepartmentRpcRequest{
			Name:                 req.Name,
			DepartmentCode:       req.DepartmentCode,
			ParentDepartmentCode: req.ParentDepartmentCode,
			OrganizationCode:     ctx.GetString("organizationCode"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var res = &model.Department{}
	_ = copier.Copy(res, departmentMessage)
	r.Success(ctx, res)
}

func (*HandlerDepartment) readDepartment(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	departmentMessage, err := rpc.DepartmentServiceClient.ReadDepartment(c,
		&department.DepartmentRpcRequest{
			DepartmentCode:   ctx.PostForm("departmentCode"),
			OrganizationCode: ctx.GetString("organizationCode"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var res = &dpt.Department{}
	_ = copier.Copy(res, departmentMessage)
	r.Success(ctx, res)
}
