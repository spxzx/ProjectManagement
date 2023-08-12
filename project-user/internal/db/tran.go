package db

// Transaction
//   - 事务的操作一定跟数据库有关 需要注入数据库的连接 gorm.db
type Transaction interface {
	Action(func(conn Conn) error) error
}
