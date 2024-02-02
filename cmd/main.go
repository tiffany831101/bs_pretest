package main

import (
	"github.com/spf13/viper"
	"github.com/tiffany831101/bs_pretest.git/config"
	"github.com/tiffany831101/bs_pretest.git/internal/database"
)

func init() {
	config.LoadConfig()

	initMongoDB()

}

func main() {

	server := StartServer()
	server.SetUpRoutes()

	server.Run()
}

func initMongoDB() {
	dbURI := viper.GetString("db.localURI")
	dbName := viper.GetString("db.name")
	database.NewDB(dbURI, dbName)
}
