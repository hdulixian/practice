package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Article struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags;"`
}

type Tag struct {
	ID       uint
	Name     string
	Articles []Article `gorm:"many2many:article_tags;"`
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
	defer func() {
		db.Migrator().DropTable(&Article{}, &Tag{}, "article_tags")
	}()
	db.AutoMigrate(&Article{}, &Tag{})

	{
		{
			db.Create(&Article{
				Title: "Python从入门到放弃",
				Tags: []Tag{
					{Name: "脚本"},
					{Name: "Python"},
				},
			})

			db.Create([]Tag{
				{Name: "Golang"},
				{Name: "Java"},
				{Name: "Database"},
				{Name: "MySQL"},
				{Name: "MongoDB"},
				{Name: "Redis"},
				{Name: "ClickHouse"},
				{Name: "Elasticsearch"},
			})
			var tags []Tag
			db.Find(&tags, "name IN ?", []string{"MySQL", "MongoDB", "Redis", "ClickHouse", "Elasticsearch"})
			db.Create(&Article{
				Title: "数据库",
				Tags:  tags,
			})
		}
	}

	{
		{
			db.Create([]Article{
				{Title: "Golang并发编程"},
				{Title: "MySQL 8.0实战"},
				{Title: "Redis开发与运维"},
				{Title: "MongoDB核心原理与实践"},
				{Title: "ClickHouse核心原理与实践"},
				{Title: "一本讲透Elasticsearch"},
			})
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"Golang"})
				var article Article
				db.Take(&article, "title = ?", "Golang并发编程")
				db.Model(&article).Association("Tags").Append(tags)
			}
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"DataBase", "MySQL"})
				var article Article
				db.Take(&article, "title = ?", "MySQL 8.0实战")
				db.Model(&article).Association("Tags").Append(tags)
			}
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"DataBase", "Redis"})
				var article Article
				db.Take(&article, "title = ?", "Redis开发与运维")
				db.Model(&article).Association("Tags").Append(tags)
			}
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"DataBase", "MongoDB"})
				var article Article
				db.Take(&article, "title = ?", "MongoDB核心原理与实践")
				db.Model(&article).Association("Tags").Append(tags)
			}
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"DataBase", "ClickHouse"})
				var article Article
				db.Take(&article, "title = ?", "ClickHouse核心原理与实践")
				db.Model(&article).Association("Tags").Append(tags)
			}
			{
				var tags []Tag
				db.Find(&tags, "name IN ?", []string{"DataBase", "Elasticsearch"})
				var article Article
				db.Take(&article, "title = ?", "一本讲透Elasticsearch")
				db.Model(&article).Association("Tags").Append(tags)
			}
		}
		{
			var tag Tag
			db.Preload("Articles").Take(&tag, "name = ?", "Database")
			fmt.Println("\n============= 1 =============")
			fmt.Println(tag)
		}
	}

	{
		{
			var article Article
			db.Preload("Tags").Take(&article, "title = ?", "数据库")
			fmt.Println("\n============= 2 =============")
			fmt.Println(article)

			var tags []Tag
			db.Find(&tags, "name IN ?", []string{"Database"})

			db.Model(&article).Association("Tags").Append(tags)
			fmt.Println("\n============= 3 =============")
			fmt.Println(article)
		}
		{
			var tag Tag
			db.Preload("Articles").Take(&tag, "name = ?", "Database")
			fmt.Println("\n============= 4 =============")
			fmt.Println(tag)
		}
	}

	{
		{
			var article Article
			db.Preload("Tags").Take(&article, "title = ?", "数据库")

			var tags []Tag
			db.Find(&tags, "name IN ?", []string{"Database"})

			db.Model(&article).Association("Tags").Replace(tags)
			fmt.Println("\n============= 5 =============")
			fmt.Println(article)
		}
		{
			var article Article
			db.Preload("Tags").Take(&article, "title = ?", "数据库")
			db.Model(&article).Association("Tags").Delete(article.Tags)
			fmt.Println("\n============= 6 =============")
			fmt.Println(article)
		}
	}

	{
		var article Article
		db.Preload("Tags").Take(&article, "title = ?", "数据库")
		fmt.Println("\n============= 7 =============")
		fmt.Println(article)

		var tag Tag
		db.Preload("Articles").Take(&tag, "name = ?", "Database")
		fmt.Println("\n============= 8 =============")
		fmt.Println(tag)
	}
}
