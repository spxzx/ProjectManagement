package repo

import (
	"context"
	"github.com/spxzx/project-project/internal/db"
	"github.com/spxzx/project-project/pkg/model"
)

type ProjectRepo interface {
	FindMyProjectListByMemId(ctx context.Context, id int64, condition string, page int64, pageSize int64) ([]*model.ProjectMemberUnion, int64, error)
	FindCollectProjectListByMemId(ctx context.Context, id int64, page int64, size int64) ([]*model.ProjectMemberUnion, int64, error)
	SaveProject(ctx context.Context, conn db.Conn, project *model.Project) error
	SaveProjectMember(ctx context.Context, conn db.Conn, pm *model.ProjectMember) error
	FindProjectByPidAndMemId(ctx context.Context, pid, mid int64) (*model.ProjectMemberUnion, error)
	FindCollectByPidAndMemId(ctx context.Context, pid, mid int64) (bool, error)
	UpdateProjectDeleted(ctx context.Context, pid int64, deleted bool) error
	SaveProjectCollected(ctx context.Context, pc *model.ProjectCollection) error
	DeleteProjectCollected(ctx context.Context, mid, pid int64) error
	UpdateProject(ctx context.Context, pj *model.Project) error
	FindProjectMemberByPid(ctx context.Context, code int64) ([]*model.ProjectMember, int64, error)
	FindProjectById(ctx context.Context, code int64) (*model.Project, error)
	FindProjectByIds(ctx context.Context, pids []int64) ([]*model.Project, error)
}
