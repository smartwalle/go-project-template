package entity

type User struct {
	Id   int64  `sql:"id,auto_increment"`
	Name string `sql:"name"`
}

func (User) TableName() string {
	return "user"
}

func (User) PrimaryKey() string {
	return "id"
}
