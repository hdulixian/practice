package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID       uint   `gorm:"size:4"`
	Name     string `gorm:"size:8"`
	Articles []Article
}

type Article struct {
	ID     uint   `gorm:"size:4"`
	Title  string `gorm:"size:16"`
	UserID uint   `gorm:"size:4"`
}

// type User struct {
// 	ID       uint      `gorm:"size:4"`
// 	Name     string    `gorm:"size:8"`
// 	Articles []Article `gorm:"foreignKey:UID;references:ID"`
// }
// type Article struct {
// 	ID    uint   `gorm:"size:4"`
// 	Title string `gorm:"size:16"`
// 	UID   uint   `gorm:"size:4"`
// }

// type User struct {
// 	ID       uint      `gorm:"size:4"`
// 	Name     string    `gorm:"unique;size:8"`
// 	Articles []Article `gorm:"foreignKey:Username;references:Name"`
// }
// type Article struct {
// 	ID       uint   `gorm:"size:4"`
// 	Title    string `gorm:"size:16"`
// 	Username string `gorm:"size:8"`
// }

var (
	db  *gorm.DB
	err error
)

func init() {
	dsn := "root:12345678@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// defer func() {
	// 	db.Migrator().DropTable(&User{}, &Article{})
	// }()

	db.AutoMigrate(&User{}, &Article{})

	{
		// user 和 article 一起创建
		db.Create(&User{
			Name: "Luna",
			Articles: []Article{
				{Title: "Python从入门到精通"},
				{Title: "Java从入门到精通"},
				{Title: "Golang从入门到精通"},
			},
		})

		// 先创建user 再创建article并关联user
		var user = User{Name: "Nick"}
		db.Create(&user)
		db.Create(&Article{
			Title:  "MySQL从删库到跑路",
			UserID: user.ID,
		})

		// 先创建article 再创建user并关联article
		db.Select("title").Create([]Article{
			{Title: "PGC2023"},
			{Title: "PGS5"},
			{Title: "PGS6"},
			{Title: "PGC2024"},
			{Title: "PGS7"},
			{Title: "PGS8"},
		})
		var articles []Article
		db.Where("title IN ?", []string{"PGC2023", "PGS5"}).Find(&articles)
		db.Create(&User{
			Name:     "Summer",
			Articles: articles,
		})

		// 创建user 再给user并添加article关联 之后再修改关联的article
		{
			var user = User{Name: "XDD"}
			db.Create(&user)

			var articles []Article
			db.Find(&articles, "title IN ?", []string{"PGS6", "PGC2024"})
			db.Model(&user).Association("Articles").Append(articles)

			db.Find(&articles, "title IN ?", []string{"PGS7", "PGS8"})
			db.Model(&user).Association("Articles").Replace(articles)
		}
	}

	{
		{
			var user User
			db.Preload("Articles").Take(&user, "name = ?", "Nick")
			fmt.Println(user)
		}
		{
			var user User
			db.Preload("Articles", "id >= ?", 2).Take(&user, 1)
			fmt.Println(user)
		}
		{
			var user User
			db.Preload("Articles", func(db *gorm.DB) *gorm.DB {
				return db.Where("id > ?", 2)
			}).Take(&user, 1)
			fmt.Println(user)
		}
	}

	{
		{
			var user User
			db.Take(&user, "name = ?", "Summer")
			// 删除user及其关联的article
			db.Select("Articles").Delete(&user)
		}
		{
			var user User
			db.Preload("Articles").Take(&user, "name = ?", "XDD")
			// 解除user与article的关联关系
			db.Model(&user).Association("Articles").Delete(&user.Articles)
			db.Delete(&user)
		}
	}
}
