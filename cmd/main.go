package main

import (
	"log"
	"microblog/internal/config"
	"microblog/internal/database"
	"microblog/internal/router"
	"microblog/internal/util"
)

func main() {
	// Загружаем конфигурацию с проверкой ошибок
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Инициализируем JWT
	util.InitJWT(cfg)

	// Инициализируем базу данных
	database.InitDB(cfg)

	// Настраиваем роутер
	r := router.Routers()

	// Используем порт из конфигурации
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s", cfg.Server.Port)

	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
