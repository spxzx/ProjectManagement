package login_service

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	common "github.com/spxzx/project-common"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-common/jwts"
	"github.com/spxzx/project-common/tms"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-user/config"
	"github.com/spxzx/project-user/internal/dao"
	"github.com/spxzx/project-user/internal/domain"
	"github.com/spxzx/project-user/internal/repo"
	"github.com/spxzx/project-user/pkg/data"
	"github.com/spxzx/project-user/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type Service struct {
	cache repo.Cache
	lcd   *domain.LoginCaptchaDomain
	ud    *domain.UserDomain
	od    *domain.OrganizationDomain
	login.UnimplementedLoginServiceServer
}

func New() *Service {
	return &Service{
		cache: dao.Rc,
		lcd:   domain.NewLoginCaptchaDomain(),
		ud:    domain.NewUserDomain(),
		od:    domain.NewOrganizationDomain(),
	}
}

func (s *Service) GetCaptcha(_ context.Context, req *login.CaptchaRequest) (*login.CaptchaResponse, error) {
	// 1.获取参数
	mobile := req.Mobile
	// 2.校验参数
	if !common.VerifyMobile(mobile) {
		return &login.CaptchaResponse{}, errs.GrpcError(data.NoLegalMobile)
	}
	// 3.生成验证码（随机4位或者6位），这里随机一下6位验证码
	//code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	code := "123456"
	// 4.调用短信平台（三方-放入go协程中执行，接口可以快速响应），这里模拟一下
	go func() {
		time.Sleep(time.Second)
		// 5.存储验证码到redis中，过期时间为15分钟
		//   假设redis后续缓存可能存在mysql当中，也可能存在mongo当中
		if err := s.lcd.PutCaptcha(data.RegisterKey+mobile, code, 15*time.Minute); err != nil {
			zap.L().Info("error in saving the verification code into redis, cause by: " + err.Error())
		}
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (s *Service) Register(_ context.Context, req *login.RegisterRequest) (*emptypb.Empty, error) {
	// 1.校验验证码
	if err := s.lcd.CheckCaptcha(req.Mobile, req.Captcha); err != nil {
		return &emptypb.Empty{}, err
	}
	// 2.校验业务逻辑（邮箱是否被注册、账号是否被注册、手机号是否被注册...）
	if err := s.ud.CheckUserExist(req.Email, req.Mobile); err != nil {
		return &emptypb.Empty{}, err
	}
	// 3.执行业务 将数据存储数据库中
	// 4.返回结果
	return &emptypb.Empty{}, s.ud.SaveMember(req)
}

func (s *Service) Login(_ context.Context, req *login.LoginRequest) (*login.LoginResponse, error) {
	// 1.校验账号和密码
	mem, err := s.ud.LoginCheck(req)
	if err != nil {
		return &login.LoginResponse{}, err
	}
	// 2.根据用户 ID 查询用户组织
	org, err := s.od.FindOrganizationByMemberId(mem.Id)
	if err != nil {
		return &login.LoginResponse{}, err
	}
	// 3.用 jwt 生成 token
	memIdStr := strconv.Itoa(int(mem.Id))
	token := jwts.CreateToken(memIdStr, config.Conf.Token, req.Ip)
	// 4.返回数据
	loginMember := &login.Member{}
	_ = copier.Copy(loginMember, mem)
	loginMember.Code, _ = encrypts.EncryptInt64(mem.Id, data.AESKey)
	loginMember.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	loginMember.CreateTime = tms.FormatByMill(mem.CreateTime)

	var loginOrg []*login.Organization
	_ = copier.Copy(&loginOrg, org)
	oMap := model.ToMap(org)
	for _, o := range loginOrg {
		o.Code, _ = encrypts.EncryptInt64(o.Id, data.AESKey)
		o.OwnerCode = loginMember.Code
		o.CreateTime = tms.FormatByMill(oMap[o.Id].CreateTime)
	}
	if len(loginOrg) > 0 {
		loginMember.OrganizationCode = loginOrg[0].Code
	}

	go func() {
		memJson, _ := json.Marshal(mem)
		orgJson, _ := json.Marshal(org)
		ct, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = s.cache.Put(ct, data.Member+"::"+memIdStr, string(memJson), config.Conf.Token.AccessExp)
		_ = s.cache.Put(ct, data.Organization+"::"+memIdStr, string(orgJson), config.Conf.Token.AccessExp)
	}()

	// 可以给 token 加一层加密处理
	return &login.LoginResponse{
		Member: loginMember,
		TokenList: &login.Token{
			AccessToken:    token.AccessToken,
			RefreshToken:   token.RefreshToken,
			TokenType:      "bearer",
			AccessTokenExp: token.AccessExp,
		},
		OrganizationList: loginOrg,
	}, nil
}

func (s *Service) TokenVerify(ctx context.Context, req *login.TokenVerifyRequest) (*login.TokenVerifyResponse, error) {
	token := req.Token
	if token == "" {
		return &login.TokenVerifyResponse{}, errs.GrpcError(data.NoLogin)
	}
	parseToken, err := jwts.ParseToken(token, config.Conf.Token.AccessSecret, req.Ip)
	if err != nil {
		zap.L().Error("tokenVerify ParseToken error, cause by: ", zap.Error(err))
		return &login.TokenVerifyResponse{}, errs.GrpcError(data.NoLogin)
	}
	_, err = strconv.ParseInt(parseToken, 10, 64)
	if err != nil {
		zap.L().Error("tokenVerify ParseInt error, cause by: ", zap.Error(err))
		return &login.TokenVerifyResponse{}, errs.GrpcError(data.NoLogin)
	}

	// Redis 缓存处理，加速响应
	mem := &login.Member{}
	memJson, err := s.cache.Get(ctx, data.Member+"::"+parseToken)
	if err != nil {
		zap.L().Error("tokenVerify redis Get Member error", zap.Error(err))
		return nil, errs.GrpcError(data.NoLogin)
	}
	_ = json.Unmarshal([]byte(memJson), &mem)
	var org []*login.Organization
	orgJson, err := s.cache.Get(ctx, data.Organization+"::"+parseToken)
	if err != nil {
		zap.L().Error("tokenVerify redis Get Member error", zap.Error(err))
		return nil, errs.GrpcError(data.NoLogin)
	}
	if memJson == "" || orgJson == "" {
		zap.L().Error("tokenVerify redis Get Member expire")
		return nil, errs.GrpcError(data.NoLogin)
	}
	_ = json.Unmarshal([]byte(orgJson), &org)
	resp := &login.Member{}
	_ = copier.Copy(resp, mem)
	if len(org) > 0 {
		resp.OrganizationCode, _ = encrypts.EncryptInt64(org[0].Id, data.AESKey)
	}
	return &login.TokenVerifyResponse{Member: resp}, nil
}

func (s *Service) GetOrgList(_ context.Context, req *login.OrgRequest) (*login.OrgResponse, error) {
	orgList, err := s.od.FindOrganizationByMemberId(req.MemberId)
	if err != nil {
		return &login.OrgResponse{}, err
	}
	var orgs []*login.Organization
	_ = copier.Copy(&orgs, orgList)
	for _, o := range orgs {
		o.Code, _ = encrypts.EncryptInt64(o.MemberId, data.AESKey)
	}
	return &login.OrgResponse{OrganizationList: orgs}, nil
}

func (s *Service) FindMemInfoById(_ context.Context, req *login.MemRequest) (*login.Member, error) {
	mem, err := s.ud.GetMemberById(req.MemberId)
	if err != nil {
		return &login.Member{}, err
	}
	resp := &login.Member{}
	_ = copier.Copy(resp, mem)
	orgs, err := s.od.FindOrganizationByMemberId(req.MemberId)
	if err != nil {
		return &login.Member{}, err
	}
	if len(orgs) > 0 {
		resp.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, data.AESKey)
	}
	resp.Code = encrypts.EncryptNoErr(mem.Id, data.AESKey)
	return resp, nil
}

func (s *Service) FindMemInfoByIds(_ context.Context, req *login.MemRequest) (*login.MemberInfoResponse, error) {
	memList, err := s.ud.FindMemberByIds(req.MemIds)
	if err != nil {
		return &login.MemberInfoResponse{}, err
	}
	var resp []*login.Member
	_ = copier.Copy(&resp, memList)
	return &login.MemberInfoResponse{List: resp}, nil
}
