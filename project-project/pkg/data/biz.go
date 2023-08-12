package data

const (
	Normal      = 1
	NonPersonal = 0
	Personal    = 1
	AESKey      = "suhvcthjkpdcvgtrfujnhgzw"
)

const (
	NotDeleted = iota
	Deleted
)
const (
	NotArchive = iota
	Archive
)
const (
	Open = iota
	Private
	Custom
)
const (
	Default = "default"
	Simple  = "simple"
)
const (
	NotCollected = iota
	Collected
)
const (
	NotOwner = iota
	Owner
)
const (
	NotExecutor = iota
	Executor
)
const (
	NotCanRead = iota
	CanRead
)
const (
	AssignTo = iota + 1
	MemberCode
	CreateBy
)
const (
	NotComment = iota
	Comment
)
