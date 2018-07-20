package user

type User struct {
	Id        int    `json:"id"          sql:"id"`
	Username  string `json:"username"    sql:"username"`
	LastName  string `json:"last_name"   sql:"last_name"`
	FirstName string `json:"first_name"  sql:"first_name"`
}

type UserService interface {
	User(id int) (*User, error)
}
