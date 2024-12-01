package user

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type Page struct {
	PageSize int `form:"pagesize" binding:"required"`
	PageNum  int `form:"pagenum" binding:"required"`
}

type GetUserListReq struct {
	Page
}

type GetUserInfoReq struct {
	Id string `uri:"id" form:"id" binding:"alphanum"`
}

type CreateUserReq struct {
	Name string `json:"name,omitempty" binding:"name"`
	Addr string `json:"addr,omitempty"`
}

type UpdateUserReq struct {
	Name string `json:"name,omitempty" binding:"name"`
	Addr string `json:"addr,omitempty"`
}
