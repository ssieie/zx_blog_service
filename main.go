package main

import (
	"blog_service/model"
	"fmt"
	"log"
)

func main() {
	err := model.InitDB()
	if err != nil {
		log.Fatalf("init DB err %s \n", err.Error())
	}

	fmt.Println("!@#!@#")
}
