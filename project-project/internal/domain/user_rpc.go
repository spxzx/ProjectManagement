package domain

import (
	"context"
	"github.com/spxzx/project-grpc/user/login"
	"github.com/spxzx/project-project/internal/rpc"
)

type UserRpcDomain struct {
	lsc login.LoginServiceClient
}

func NewUserRpcDomain() *UserRpcDomain {
	return &UserRpcDomain{lsc: rpc.LoginServiceClient}
}

func (u *UserRpcDomain) FindMemInfoById(ctx context.Context, mid int64) (*login.Member, error) {
	//mem, err := u.lsc.FindMemInfoById(ctx, req)
	mem, err := rpc.LoginServiceClient.FindMemInfoById(ctx, &login.MemRequest{MemberId: mid})
	if err != nil {
		return nil, err
	}
	return mem, nil
}

func (u *UserRpcDomain) FindMemInfoByIds(ctx context.Context, req *login.MemRequest) (*login.MemberInfoResponse,
	map[int64]*login.Member, error) {
	//memList, err := u.lsc.FindMemInfoByIds(ctx, req)
	memList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	mMap := make(map[int64]*login.Member)
	for _, v := range memList.List {
		mMap[v.Id] = v
	}
	return memList, mMap, nil
}
