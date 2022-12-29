package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pantabank/t-shirts-shop/configs"
	"github.com/pantabank/t-shirts-shop/internals/servers"
	databases "github.com/pantabank/t-shirts-shop/pkg/databases"
)

func main() {
	// Env can prase as a flag args for dyanmics value
	// flag.Parse()
	// defaultPort := "3000"

	// Load dotenv config
	// This can be dynamically with .env using flag
	// defaultEnvPath = ../.env
	// envPath := flag.String("env_path", defaultEnvPath, "The server's .env path")

	if err := godotenv.Load("../.env"); err != nil {
		panic(err.Error())
	}

	// Config like this is more cleaner
	// cfg := &configs.Configs{
	// 	PostgreSQL: configs.PostgreSQL{
	// 		Host:     os.Getenv("DB_HOST"),
	// 		Port:     os.Getenv("DB_PORT"),
	// 		Protocol: os.Getenv("DB_PROTOCOL"),
	// 		Username: os.Getenv("DB_USERNAME"),
	// 		Password: os.Getenv("DB_PASSWORD"),
	// 		Database: os.Getenv("DB_DATABASE"),
	// 	},
	// 	App: configs.Fiber{
	// 		Host: os.Getenv("FIBER_HOST"),
	// 		Port: flag.String("port", defaultPort, "The server's port"), // Must change to pointer if using a flag
	// 	},
	// }

	cfg := new(configs.Configs)
	// Fiber configs
	cfg.App.Host = os.Getenv("FIBER_HOST")
	cfg.App.Port = os.Getenv("FIBER_PORT")

	// // Database Configs
	cfg.PostgreSQL.Host = os.Getenv("DB_HOST")
	cfg.PostgreSQL.Port = os.Getenv("DB_PORT")
	cfg.PostgreSQL.Protocol = os.Getenv("DB_PROTOCOL")
	cfg.PostgreSQL.Username = os.Getenv("DB_USERNAME")
	cfg.PostgreSQL.Password = os.Getenv("DB_PASSWORD")
	cfg.PostgreSQL.Database = os.Getenv("DB_DATABASE")

	// New Database
	db, err := databases.NewPostgreSQLDBConnection(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	s := servers.NewServer(cfg, db)
	s.Start()
}
