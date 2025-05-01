package main

import (
	"fmt"
	"time"

	"investbot/pkg/services"
)

type User struct {
	ID   int
	Name string
}

func main() {
	cache, _ := services.NewBadgerCacheService()

	user := User{ID: 1, Name: "Alice"}
	cache.Set("user:1", user, time.Minute)

	var u User
	err := cache.Get("user:1", &u)
	if err == nil {
		fmt.Println(u.Name) // Output: Alice
	}

	time.Sleep(time.Second * 65)

	err = cache.Get("user:1", &u)
	if err == nil {
		fmt.Println(u.Name) // Output: Alice
	}
}
