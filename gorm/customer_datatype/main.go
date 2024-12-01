package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Info struct {
	Status string `json:"status"`
	Addr   string `json:"addr"`
	Age    int    `json:"age"`
}

func (info Info) Value() (driver.Value, error) {
	return json.Marshal(info)
}

func (info *Info) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}
	return json.Unmarshal(bytes, info)
}

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Info Info   `json:"info" gorm:"type:string"`
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
		db.Migrator().DropTable(&User{}, &Task{})
	}()

	{
		db.AutoMigrate(&User{})

		db.Create(&User{
			Name: "Luna",
			Info: Info{
				Status: "Success",
				Addr:   "Beijing",
				Age:    23,
			},
		})

		var user User
		db.Take(&user)
		fmt.Println(user)

		userStr, _ := json.Marshal(user)
		fmt.Println(string(userStr))
	}

	{
		db.AutoMigrate(&Task{})

		db.Create(&Task{
			Name:   "sca-scan",
			Status: Running,
		})

		var task Task
		db.Find(&task)
		fmt.Println(task)

		taskStr, _ := json.Marshal(task)
		fmt.Println(string(taskStr))
	}
}

type Status int

const (
	Running Status = iota + 1
	Successed
	Failed
)

func (s Status) MarshalJSON() ([]byte, error) {
	var status string
	switch s {
	case Running:
		status = "Running"
	case Successed:
		status = "Successed"
	case Failed:
		status = "Failed"
	}
	return json.Marshal(status)
}

type Task struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status" gorm:"type:tinyint"`
}
