package main

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

type Student struct {
	ID     uint   `gorm:"size:3"`
	Name   string `gorm:"size:8"`
	Age    int    `gorm:"size:3"`
	Gender bool
	Email  *string `gorm:"size:32"`
}

func main() {
	defer func() {
		db.Migrator().DropTable(&Student{})
	}()

	db.AutoMigrate(&Student{})

	{
		var stus []Student
		for i := 1; i <= 10; i++ {
			stus = append(stus, Student{
				Name:   strconv.Itoa(i),
				Age:    10 + i,
				Gender: true,
				Email:  nil,
			})
		}
		db.Debug().CreateInBatches(stus, 10)

		var stu Student
		// Normal Query
		db.Debug().Take(&stu, "name = ?", "2' OR name = '1")
		db.Raw("SELECT * FROM `students` WHERE name = ? LIMIT 1", "2' OR name = '1").Take(&stu)

		// SQL Injection
		db.Raw(fmt.Sprintf("SELECT * FROM `students` WHERE name = '%s' LIMIT 1", "2' OR name = '1")).Take(&stu)
		fmt.Println(stu)

		{
			stu := Student{
				ID:   3,
				Name: "4",
			}
			// 结构体查询只使用主键
			// SELECT * FROM `students` WHERE `students`.`id` = 3 LIMIT 1
			db.Debug().Take(&stu)
			fmt.Println(stu)
		}

		{
			var stus []Student
			// where 和 take 查询条件聚合
			// SELECT * FROM `students` WHERE name IN ('2','3') AND `students`.`id` IN (1,2,3)
			db.Debug().Where("name IN ?", []string{"2", "3"}).Find(&stus, 1, 2, 3)
			fmt.Println(stus)
		}
	}

	{
		stus := []Student{
			{Name: "李元芳", Age: 32, Email: PtrString("lyf@yf.com"), Gender: true},
			{Name: "张武", Age: 18, Email: PtrString("zhangwu@lly.cn"), Gender: true},
			{Name: "枫枫", Age: 23, Email: PtrString("ff@yahoo.com"), Gender: true},
			{Name: "刘大", Age: 54, Email: PtrString("liuda@qq.com"), Gender: true},
			{Name: "李武", Age: 23, Email: PtrString("liwu@lly.cn"), Gender: true},
			{Name: "李琦", Age: 14, Email: PtrString("liqi@lly.cn"), Gender: false},
			{Name: "晓梅", Age: 25, Email: PtrString("xiaomeo@sl.com"), Gender: false},
			{Name: "如燕", Age: 26, Email: PtrString("ruyan@yf.com"), Gender: false},
			{Name: "魔灵", Age: 21, Email: PtrString("moling@sl.com"), Gender: true},
		}
		db.Create(&stus)

		{
			var stu Student
			// SELECT * FROM `students` WHERE `students`.`name` = '李元芳' AND `students`.`age` = 32 LIMIT 1
			db.Debug().Where(&Student{Name: "李元芳", Age: 32}).Take(&stu)
		}
		{
			{
				var stu Student
				// 结构体查询时零值字段会被自动过滤掉
				// SELECT * FROM `students` WHERE `students`.`name` = '李元芳' LIMIT 1
				db.Debug().Where(&Student{Name: "李元芳", Age: 0}).Take(&stu)
			}
			{
				var stu Student
				// 结构体查询时零值字段会被自动过滤掉
				// SELECT * FROM `students` WHERE `students`.`name` = '李元芳' LIMIT 1
				db.Debug().Take(&stu, &Student{Name: "李元芳", Age: 0})
			}
		}

		{
			type genderAggs struct {
				Gender bool
				Count  int
				AvgAge float32
			}
			var aggs []genderAggs
			db.Debug().Model(&Student{}).Select("gender", "COUNT(1) AS count", "AVG(age) AS avg_age").Group("gender").Scan(&aggs)
			fmt.Println(aggs)
		}

		{
			var stus []Student
			db.Debug().Where("age > (?)", db.Model(&Student{}).Select("AVG(age)")).Find(&stus)
			fmt.Println(stus)
		}

		{
			var stus []Student
			db.Debug().Scopes(AgeScope(17, 30)).Find(&stus)
			fmt.Println(stus)
		}
	}
}

func PtrString(email string) *string {
	return &email
}

func AgeScope(minAge, maxAge int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age BETWEEN ? AND ?", minAge, maxAge)
	}
}
