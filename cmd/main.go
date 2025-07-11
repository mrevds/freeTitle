package main

import (
	"microblog/internal/config"
	"microblog/internal/database"
	"microblog/internal/router"
	"microblog/internal/util"
)

func main() {
	cfg := config.LoadConfig()

	util.InitJWT(cfg)

	database.InitDB(cfg)

	r := router.Routers()
	if err := r.Run(":8080"); err != nil {
		panic("ошибка запуска сервера: " + err.Error())
	}
}
