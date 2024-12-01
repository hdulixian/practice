package main

import (
	"io"
	"net/http"
	"os"

	"gin/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	govalidator, _ := binding.Validator.Engine().(*validator.Validate)
	govalidator.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		name := fl.Field().Interface().(string)
		return len(name) >= 3
	})
}

func main() {
	gin.ForceConsoleColor()

	file, _ := os.OpenFile("gin.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)

	router := gin.Default()

	router.StaticFS("/static", http.Dir("/Users/emmcd/Desktop/3dparty"))

	{
		router.LoadHTMLGlob("html/**/*")
		router.GET("/welcome", func(c *gin.Context) {
			c.HTML(http.StatusOK, "welcome.html", nil)
		})
		router.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{"detail": "Gin框架"})
		})
	}

	{
		router.GET("/go.work", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/static/lru.go")
			// c.Redirect(http.StatusFound, "http://www.baidu.com")
		})
	}

	{
		router.GET("/go.work.sum", func(c *gin.Context) {
			c.Header("Content-Type", "application/octet-stream")
			c.File("../go.work.sum")
		})
	}

	{
		// MaxMultipartMemory 用于限制处理 multipart 表单时可以分配的最大内存
		// 如果上传的文件大小超过了这个限制，超出部分将会被存储在临时文件中，而不是内存里
		// 这个设置并不直接限制上传文件的大小，而是限制了在内存中可以存储 multipart 表单数据的大小
		router.MaxMultipartMemory = 2 << 20

		router.POST("/upload", func(c *gin.Context) {
			c.Copy()
			file, _ := c.FormFile("file")
			c.SaveUploadedFile(file, "./"+file.Filename)
			c.JSON(200, gin.H{"code": 200, "msg": "上传成功"})
		})

		router.POST("/uploads", func(c *gin.Context) {
			form, _ := c.MultipartForm()
			files := form.File["files"]
			for _, file := range files {
				c.SaveUploadedFile(file, "./"+file.Filename)
			}
			c.JSON(200, gin.H{"code": 200, "msg": "上传成功"})
		})
	}

	{
		var userController *user.UserController
		{
			router.GET("/user", userController.GetUserList)
			router.GET("/user/:id", userController.GetUserInfo)
			router.POST("/user", userController.CreateUser)
			router.PATCH("/user/:id", userController.UpdateUser)
			router.DELETE("/user/:id", userController.DeleteUser)
		}
	}

	router.Run(":80")
}
