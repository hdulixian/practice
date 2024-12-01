package user

import (
	"gin/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct{}

func (*UserController) GetUserList(c *gin.Context) {
	var req GetUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, 4004, err.Error())
		return
	}
	response.Success(c, 0, "success", []User{
		{Id: uuid.NewString(), Name: "Luna", Addr: "Pairs"},
		{Id: uuid.NewString(), Name: "Nick", Addr: "London"},
		{Id: uuid.NewString(), Name: "Mike", Addr: "New York"},
	})
}

func (*UserController) GetUserInfo(c *gin.Context) {
	var req GetUserInfoReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.Fail(c, 4004, err.Error())
		return
	}
	response.Success(c, 0, "success", &User{Id: req.Id, Name: "Luna", Addr: "Pairs"})
}

func (*UserController) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 4004, err.Error())
		return
	}
	response.Success(c, 0, "success", &User{Id: uuid.NewString(), Name: req.Name, Addr: req.Addr})
}

func (*UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 4004, err.Error())
		return
	}
	response.Success(c, 0, "success", &User{Id: id, Name: req.Name, Addr: req.Addr})
}

func (*UserController) DeleteUser(c *gin.Context) {
	_ = c.Param("id")
	response.Success(c, 0, "success", nil)
}
