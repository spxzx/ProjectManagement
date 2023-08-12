package repo

import (
	"context"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/pkg/model"
)

type MemberRepo interface {
	ExistMemberByEmail(ctx context.Context, email string) (bool, error)
	ExistMemberByAccount(ctx context.Context, account string) (bool, error)
	ExistMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(ctx context.Context, conn db.Conn, mem *model.Member) error
	FindMember(ctx context.Context, account, password string) (model.Member, error)
	GetMemberById(ctx context.Context, id int64) (model.Member, error)
	FindMemberByIds(ctx context.Context, ids []int64) ([]*model.Member, error)
}
