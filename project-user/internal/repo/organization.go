package repo

import (
	"context"
	"github.com/spxzx/project-user/internal/db"
	"github.com/spxzx/project-user/pkg/model"
)

type OrganizationRepo interface {
	SaveOrganization(ctx context.Context, conn db.Conn, org *model.Organization) error
	FindOrganizationByMemberId(ctx context.Context, id int64) ([]model.Organization, error)
}
