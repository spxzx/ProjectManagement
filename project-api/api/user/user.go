package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spxzx/project-api/api/rpc"
	"github.com/spxzx/project-api/pkg/model/member"
	"github.com/spxzx/project-api/pkg/model/organization"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/r"
	"github.com/spxzx/project-grpc/user/login"
	"net/http"
	"time"
)

type HandlerUser struct{}

func New() *HandlerUser {
	return &HandlerUser{}
}

func (*HandlerUser) getCaptcha(ctx *gin.Context) {
	mobile := ctx.PostForm("mobile")
	// 下面传入这个 gin 的 ctx 就只是当作一个 nil 的 context 去传， 直接传 nil 会报错
	resp, err := rpc.LoginServiceClient.GetCaptcha(ctx, &login.CaptchaRequest{Mobile: mobile})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	r.Success(ctx, resp.Code)
}

func (*HandlerUser) register(ctx *gin.Context) {
	// 1.接收参数 - 参数模型
	var req member.RegisterReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	// 2.校验参数 - 判断参数是否合法
	if err := req.Verify(); err != nil {
		r.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}
	// 3.调用user grpc服务 - 获取响应
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &login.RegisterRequest{}
	_ = copier.Copy(in, req) // 一般不会出错，这里忽略这个错误处理
	if _, err := rpc.LoginServiceClient.Register(c, in); err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	// 4.返回结果
	r.Success(ctx)
}

func (*HandlerUser) login(ctx *gin.Context) {
	// 1.接收参数
	var req member.LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		r.Fail(ctx, http.StatusBadRequest, "参数格式有误")
		return
	}
	req.Password = encrypts.Md5(req.Password)
	// 2.调用user grpc服务
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	in := &login.LoginRequest{}
	_ = copier.Copy(in, req)
	in.Ip = GetIp(ctx)
	loginResp, err := rpc.LoginServiceClient.Login(c, in)
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	resp := &member.LoginResp{}
	_ = copier.Copy(resp, loginResp)
	r.Success(ctx, resp)
}

func (*HandlerUser) getMyOrgList(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	orgResp, err := rpc.LoginServiceClient.GetOrgList(c, &login.OrgRequest{MemberId: ctx.GetInt64("memberId")})
	if err != nil {
		e := errs.ParseGrpcError(err)
		r.Fail(ctx, int(e.Code), e.Msg)
		return
	}
	var resp []*organization.Organization
	_ = copier.Copy(&resp, orgResp.OrganizationList)
	r.Success(ctx, resp)
}

func GetIp(c *gin.Context) (ip string) {
	ip = c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return
}
