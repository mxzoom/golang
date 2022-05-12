package main

import (
	"os"
	"parser/config"
)

//
//go run main.go -port=4040 -db_url="postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=enable" -some_app_id="ID45695199" -some_url="https://ya.ru"
var cfgStr config.MyConf

func main() {
	os.Setenv("port", "8080")
	os.Setenv("db_url", "postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=disable")
	os.Setenv("some_app_id", "some_id")
	os.Setenv("some_url", "http://sentry:9000")

	cfgStr.MainParse()
}
