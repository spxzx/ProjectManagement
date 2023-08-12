package repo

import (
	"context"
	"github.com/spxzx/project-project/pkg/model"
)

type MemberAccountRepo interface {
	FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) ([]*model.MemberAccount, int64, error)
	FindByMemberId(ctx context.Context, memId int64) (*model.MemberAccount, error)
}
