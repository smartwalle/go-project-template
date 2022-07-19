package http

type UserRsp struct {
	Id        int64  `json:"id"`         // id
	Username  string `json:"username"`   // 用户名
	LastName  string `json:"last_name"`  // 姓
	FirstName string `json:"first_name"` // 名
}

type AddUserReq struct {
	Username  string `form:"username"`   // 用户名
	LastName  string `form:"last_name"`  // 姓
	FirstName string `form:"first_name"` // 名
}
