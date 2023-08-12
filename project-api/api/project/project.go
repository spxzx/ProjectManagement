package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model"
	"github.com/spxzx/project-api/pkg/model/menu"
	"github.com/spxzx/project-api/pkg/model/poj"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/project/project"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"strconv"
	"time"
)

type HandlerProject struct{}

func NewProject() *HandlerProject {
	return &HandlerProject{}
}

func (*HandlerProject) index(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	indexResp, err := rpc.ProjectServiceClient.Index(c, &emptypb.Empty{})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	resp := &menu.IndexResp{}
	_ = copier.Copy(resp, indexResp)
	r.Success(ctx, resp.Menus)
}

func (*HandlerProject) getProjectList(ctx *gin.Context) {
	page := &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	selectBy := ctx.PostForm("selectBy")
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pojResp, err := rpc.ProjectServiceClient.GetProjectList(c,
		&project.ProjectRpcRequest{
			MemberId:   ctx.GetInt64("memberId"),
			MemberName: ctx.GetString("memberName"),
			Page:       page.Page,
			PageSize:   page.PageSize,
			SelectBy:   selectBy,
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp []*poj.ProjectMemberUnion
	_ = copier.Copy(&resp, pojResp.Pm)
	if len(resp) == 0 {
		resp = []*poj.ProjectMemberUnion{}
	}
	r.Success(ctx, gin.H{
		"list":  resp,
		"total": pojResp.Total,
	})
}

func (*HandlerProject) getProjectTemplates(ctx *gin.Context) {
	page := &model.Page{}
	_ = page.Bind(ctx)
	viewType, _ := strconv.ParseInt(ctx.PostForm("viewType"), 10, 64)
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ptResp, err := rpc.ProjectServiceClient.GetProjectTemplates(c, &project.ProjectRpcRequest{
		MemberId:         ctx.GetInt64("memberId"),
		MemberName:       ctx.GetString("memberName"),
		Page:             page.Page,
		PageSize:         page.PageSize,
		OrganizationCode: ctx.GetString("organizationCode"),
		ViewType:         int32(viewType),
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var pts []*poj.ProjectTemplate
	_ = copier.Copy(&pts, ptResp.Pt)
	if len(pts) == 0 {
		pts = []*poj.ProjectTemplate{}
	}
	r.Success(ctx, gin.H{
		"list":  pts,
		"total": ptResp.Total,
	})
}

func (*HandlerProject) saveProject(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &poj.SaveProjectReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	saveResp, err := rpc.ProjectServiceClient.SaveProject(
		c, &project.ProjectRpcRequest{
			MemberId:         ctx.GetInt64("memberId"),
			OrganizationCode: ctx.GetString("organizationCode"),
			Name:             req.Name,
			TemplateCode:     req.TemplateCode,
			Description:      req.Description,
			Id:               req.Id,
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp *poj.SaveProjectResp
	copier.Copy(&resp, saveResp)
	r.Success(ctx, resp)
}

func (*HandlerProject) readProject(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	readResp, err := rpc.ProjectServiceClient.ReadProject(c,
		&project.ProjectRpcRequest{
			ProjectCode: ctx.PostForm("projectCode"),
			MemberId:    ctx.GetInt64("memberId"),
		})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	resp := &poj.ReadProjectResp{}
	_ = copier.Copy(resp, readResp)
	r.Success(ctx, resp)
}

func (*HandlerProject) recycleProject(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := rpc.ProjectServiceClient.UpdateProjectDeleted(c,
		&project.ProjectRpcRequest{
			ProjectCode: ctx.PostForm("projectCode"),
			Deleted:     true,
		},
	); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

func (*HandlerProject) recoveryProject(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := rpc.ProjectServiceClient.UpdateProjectDeleted(c,
		&project.ProjectRpcRequest{
			ProjectCode: ctx.PostForm("projectCode"),
			Deleted:     false,
		},
	); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

func (*HandlerProject) collectProject(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := rpc.ProjectServiceClient.UpdateProjectCollected(c,
		&project.ProjectRpcRequest{
			MemberId:    ctx.GetInt64("memberId"),
			ProjectCode: ctx.PostForm("projectCode"),
			CollectType: ctx.PostForm("type"),
		},
	); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

// TODO ** 这个模块肯定是存在问题的，更改为 0 的时候 gorm会默认忽略字段 **
func (*HandlerProject) editProject(ctx *gin.Context) {
	var req *poj.EditProjectReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	editReq := &project.EditProjectRequest{}
	_ = copier.Copy(editReq, req)
	editReq.MemberId = ctx.GetInt64("memberId")
	if _, err := rpc.ProjectServiceClient.EditProject(c, editReq); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx)
}

func (*HandlerProject) getLogBySelfProject(ctx *gin.Context) {
	page := &model.Page{}
	if err := page.Bind(ctx); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	logResp, err := rpc.ProjectServiceClient.GetLogBySelfProject(c, &project.ProjectRpcRequest{
		MemberId: ctx.GetInt64("memberId"),
		Page:     page.Page,
		PageSize: page.PageSize,
	})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*poj.ProjectLog
	_ = copier.Copy(&list, logResp.List)
	if list == nil || len(list) <= 0 {
		list = []*poj.ProjectLog{}
	}
	r.Success(ctx, list)
}

func (*HandlerProject) getNodeList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := rpc.ProjectServiceClient.GetNodeList(c, &project.ProjectRpcRequest{})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var list []*poj.ProjectNodeTree
	_ = copier.Copy(&list, response.Nodes)
	r.Success(ctx, gin.H{"nodes": list})
}
