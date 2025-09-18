package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/cobrich/recommendo/config"
	"github.com/cobrich/recommendo/handlers"
	"github.com/cobrich/recommendo/repo"
	"github.com/cobrich/recommendo/router"
	"github.com/cobrich/recommendo/service"

	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Get config
	cfg := config.GetConfig()

	// Create connector
	connector := stdlib.GetConnector(*cfg)

	// Open db with connector that saves instructions for how creating
	db := sql.OpenDB(connector)

	// Close db for free space in memory
	defer db.Close()

	// Check connecting to db
	err := db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	fmt.Println("Successfully connected to database!")

	// Repos
	userRepo := repo.NewUserRepo(db)
	friendhipRepo := repo.NewFollowRepo(db)

	// Services
	userService := service.NewUserService(userRepo)
	friendshipService := service.NewFriendshipService(friendhipRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	friendshipHandler := handlers.NewFriendshiphandler(friendshipService)

	fmt.Println("Сервер запущен на http://localhost:8080")

	// Create router and set
	router := router.NewRouter(userHandler, friendshipHandler)

	// Run server in port 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}
