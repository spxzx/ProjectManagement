package domain

import (
	"context"
	"fmt"
	"github.com/spxzx/project-common/encrypts"
	"github.com/spxzx/project-common/errs"
	"github.com/spxzx/project-project/internal/dao"
	"github.com/spxzx/project-project/internal/repo"
	"github.com/spxzx/project-project/pkg/data"
	"github.com/spxzx/project-project/pkg/model"
	"go.uber.org/zap"
	"time"
)

type MemberAccountDomain struct {
	ma      repo.MemberAccountRepo
	dd      *DepartmentDomain
	userRpc *UserRpcDomain
}

func NewMemberAccountDomain() *MemberAccountDomain {
	return &MemberAccountDomain{
		ma:      dao.NewMemberAccountDao(),
		dd:      NewDepartmentDomain(),
		userRpc: NewUserRpcDomain(),
	}
}

func (m *MemberAccountDomain) GetAccountList(ctx context.Context, organizationCode string, memberId int64,
	page int64, pageSize int64, departmentCode string, searchType int32,
) ([]*model.MemberAccountDisplay, int64, error) {
	condition := ""
	organizationCodeId := encrypts.DecryptNoErr(organizationCode, data.AESKey)
	departmentCodeId := encrypts.DecryptNoErr(departmentCode, data.AESKey)
	switch searchType {
	case 1:
		condition = "status = 1"
	case 2:
		condition = "department_code = NULL"
	case 3:
		condition = "status = 0"
	case 4:
		condition = fmt.Sprintf("status = 1 and department_code = %d", departmentCodeId)
	default:
		condition = "status = 1"
	}
	list, total, err := m.ma.FindList(ctx, condition, organizationCodeId, departmentCodeId, page, pageSize)
	if err != nil {
		zap.L().Error("account domain FindList error, cause by: ", zap.Error(err))
		return nil, 0, errs.GrpcError(data.DBError)
	}
	var dList []*model.MemberAccountDisplay
	for _, v := range list {
		display := v.ToDisplay()
		memberInfo, _ := m.userRpc.FindMemInfoById(ctx, v.MemberCode)
		display.Avatar = memberInfo.Avatar
		if v.DepartmentCode > 0 {
			department, err := m.dd.FindDepartmentById(ctx, v.DepartmentCode)
			if err != nil {
				zap.L().Error("account domain FindDepartmentById error, cause by: ", zap.Error(err))
				return nil, 0, errs.GrpcError(data.DBError)
			}
			display.Departments = department.Name
		}
		dList = append(dList, display)
	}
	return dList, total, nil
}

func (m *MemberAccountDomain) FindAccount(memId int64) (*model.MemberAccount, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberAccount, err := m.ma.FindByMemberId(c, memId)
	if err != nil {
		return nil, errs.GrpcError(data.DBError)
	}
	return memberAccount, nil
}
