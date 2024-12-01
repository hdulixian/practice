package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       uint
	Name     string
	Age      int
	Gender   bool
	UserInfo UserInfo // `gorm:"foreignKey:UserID;references:ID"`
}

type UserInfo struct {
	ID     uint
	UserID uint
	Addr   string
	Like   string
}

var (
	db  *gorm.DB
	err error
)

func init() {
	dsn := "root:12345678@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// defer func() {
	// 	db.Migrator().DropTable(&User{}, &UserInfo{})
	// }()

	db.AutoMigrate(&User{}, &UserInfo{})

	{
		db.Create(&User{
			Name:   "Luna",
			Age:    23,
			Gender: false,
			UserInfo: UserInfo{
				Addr: "Beijing",
				Like: "Tennis",
			},
		})
	}

	{
		{
			var user = User{
				Name:   "Summer",
				Age:    24,
				Gender: true,
			}
			db.Create(&user)
			db.Create(&UserInfo{
				UserID: user.ID,
				Addr:   "Shenzhen",
				Like:   "Hike",
			})
		}

		{
			var user = User{
				Name:   "Nick",
				Age:    27,
				Gender: true,
			}
			db.Create(&user)
			db.Debug().Model(&user).Association("UserInfo").Append(&UserInfo{
				Addr: "Shanghai",
				Like: "Swim",
			})
		}
	}

	{
		var user User
		db.Preload("UserInfo").Take(&user, "name = ?", "Summer")
		fmt.Println(user)
	}
}
