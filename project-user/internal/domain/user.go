package domain

import (
	"context"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-user/internal/dao"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/internal/repo"
	"github.com/spxzx/project-user/pkg/data"
	"github.com/spxzx/project-user/pkg/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type UserDomain struct {
	member    repo.MemberRepo
	tsc       db.Transaction
	orgDomain *OrganizationDomain
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		member:    dao.NewMemberDao(),
		tsc:       dao.NewTransaction(),
		orgDomain: NewOrganizationDomain(),
	}
}

func (u *UserDomain) CheckUserExist(email, mobile string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	exist, err := u.member.ExistMemberByEmail(ctx, email)
	if err != nil {
		zap.L().Error("register db ExistMemberByEmail error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	if exist {
		return errs.GrpcError(data.EmailExist)
	}
	exist, err = u.member.ExistMemberByAccount(ctx, email)
	if err != nil {
		zap.L().Error("register db ExistMemberByAccount error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	if exist {
		return errs.GrpcError(data.AccountExist)
	}
	exist, err = u.member.ExistMemberByMobile(ctx, mobile)
	if err != nil {
		zap.L().Error("register db ExistMemberByMobile error, cause by: ", zap.Error(err))
		return errs.GrpcError(data.DBError)
	}
	if exist {
		return errs.GrpcError(data.MobileExist)
	}
	return nil
}

func (u *UserDomain) SaveMember(req *login.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return u.tsc.Action(func(conn db.Conn) error {
		mem := &model.Member{
			Account:       req.Name,
			Password:      encrypts.Md5(req.Password),
			Name:          req.Name,
			Mobile:        req.Mobile,
			Email:         req.Email,
			CreateTime:    time.Now().UnixMilli(),
			LastLoginTime: time.Now().UnixMilli(),
			Status:        data.Normal,
		}
		if err := u.member.SaveMember(ctx, conn, mem); err != nil {
			zap.L().Error("register db SaveMember error, cause by: ", zap.Error(err))
			return errs.GrpcError(data.DBError)
		}
		if err := u.orgDomain.SaveOrganization(req, mem); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserDomain) LoginCheck(req *login.LoginRequest) (model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mem, err := u.member.FindMember(ctx, req.Account, req.Password)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return mem, errs.GrpcError(data.AccountOrPasswordError)
		}
		zap.L().Error("login db FindMember error, cause by: ", zap.Error(err))
		return mem, err
	}
	return mem, nil
}

func (u *UserDomain) GetMemberById(id int64) (model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mem, err := u.member.GetMemberById(ctx, id)
	if err != nil {
		zap.L().Error("findMemInfoById db FindMemInfoById error, cause by: ", zap.Error(err))
		return mem, errs.GrpcError(data.DBError)
	}
	return mem, nil
}

func (u *UserDomain) FindMemberByIds(ids []int64) ([]*model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	memList, err := u.member.FindMemberByIds(ctx, ids)
	if err != nil {
		zap.L().Error("findMemInfoByIds db FindMemberByIds error, cause by: ", zap.Error(err))
		return memList, errs.GrpcError(data.DBError)
	}
	return memList, nil
}
