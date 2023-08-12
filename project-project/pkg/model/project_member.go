package model

type ProjectMember struct {
	Id          int64
	ProjectCode int64
	MemberCode  int64
	JoinTime    int64
	IsOwner     int
	Authorize   string
}

func ToPMMemIdsAndMap(pms []*ProjectMember) (ids []int64, m map[int64]*ProjectMember) {
	m = make(map[int64]*ProjectMember)
	for _, v := range pms {
		ids = append(ids, v.MemberCode)
		m[v.MemberCode] = v
	}
	return
}

//type ProjectMemberInfo struct {
//	ProjectCode int64
//	MemberCode  int64
//	Name        string
//	Avatar      string
//	IsOwner     int64
//	Email       string
//}
//
//func ToPMInfoMap(pm []*ProjectMemberInfo) map[int64]*ProjectMemberInfo {
//	m := make(map[int64]*ProjectMemberInfo)
//	for _, v := range pm {
//		m[v.MemberCode] = v
//	}
//	return m
//}
