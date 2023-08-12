package model

type Organization struct {
	Id          int64
	Name        string
	Avatar      string
	Description string
	MemberId    int64
	CreateTime  int64
	Personal    int32
	Address     string
	Province    int32
	City        int32
	Area        int32
}

func ToMap(orgs []Organization) map[int64]Organization {
	m := make(map[int64]Organization)
	for _, v := range orgs {
		m[v.Id] = v
	}
	return m
}
