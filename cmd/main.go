package main

import (
	"microblog/internal/database"
	"microblog/internal/router"
)

func main() {
	database.InitDB()
	r := router.Routers()
	r.Run(":8080")
}
