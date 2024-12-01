package main

import (
	"log"
	"os"
)

type User struct {
	Name string
	Age  int
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.SetPrefix("[GIN] ")
	logger.SetFlags(log.LstdFlags | log.Lshortfile)

	u := User{
		Name: "dj",
		Age:  18,
	}
	logger.Printf("%s login, age:%d", u.Name, u.Age)
}
