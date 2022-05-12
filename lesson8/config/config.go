package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

/*
test args here...


port: 8080
db_url: postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=disable
some_app_id: testid
some_url: http://sentry:9000
*/

type MyConf struct {
	Port        *int
	Db_url      *string
	Some_app_id *string
	Some_url    *string
	Checked     bool
	Err         []error
	Source      string
}

func (cfg *MyConf) ParseFlags() {

	cfg.Port = flag.Int("port", 0, "insert port number here")
	cfg.Db_url = flag.String("db_url", "", "insert db_url here")
	cfg.Some_app_id = flag.String("some_app_id", "", "some app_id_here")
	cfg.Some_url = flag.String("some_url", "", "insert some_url here")
	flag.Parse()
	cfg.Source = "flag"
}

func (cfg *MyConf) ParseEnv() {
	tPort, tPort_err := strconv.Atoi((os.Getenv("port")))
	if tPort_err != nil {
		fmt.Println("Переменная окружения - `port` задана некорректно ")
	}

	tDb_url := os.Getenv("db_url")
	tSome_app_id := os.Getenv("some_app_id")
	tSome_url := os.Getenv("some_url")
	cfg.Port = &tPort
	cfg.Db_url = &tDb_url
	cfg.Some_app_id = &tSome_app_id
	cfg.Some_url = &tSome_url
	cfg.Source = "env"
}

func (cfg *MyConf) CheckArgs() {
	if len(*cfg.Some_app_id) == 0 {
		cfg.Err = append(cfg.Err, errors.New("введите идентификатор приложения -some_app_id=[a-Z, 0-9]"))
	}
	if *cfg.Port == 0 {
		cfg.Err = append(cfg.Err, errors.New("введите номер порта -port=[0-9]{1-6}"))
	}
	_, ok := url.ParseRequestURI(*cfg.Some_url)
	if ok != nil {
		cfg.Err = append(cfg.Err, errors.New("введите валидный url -some_url=[url]"))
	}
	if len(cfg.Err) == 0 {
		cfg.Checked = true
	}
	re, _ := regexp.Compile(`^postgres://.+:.+@.+:\d{0,5}/.+\?sslmode=(enable$|disable$)`)
	if !re.MatchString(*cfg.Db_url) {
		cfg.Err = append(cfg.Err, errors.New("введите валидную строку подключения к БД, например [postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=disable]"))
	}
	if len(cfg.Err) == 0 {
		cfg.Checked = true
	} else {
		cfg.Checked = false
	}
}

func (cfg *MyConf) PrintArgs() {
	if cfg.Checked {
		fmt.Printf("%d\n%s\n%s\n%s\n", *cfg.Port, *cfg.Db_url, *cfg.Some_app_id, *cfg.Some_url)
		if cfg.Source == "env" {
			fmt.Println("Источник данных: переменные окружения")
		} else {
			fmt.Println("Источник данных: аргументы командной строки")
		}
	} else {
		for i := 0; i < len(cfg.Err); i++ {
			fmt.Println(cfg.Err[i])

		}
	}
}

func (cfg *MyConf) MainParse() {
	if len(os.Args[:]) > 1 {
		cfg.ParseFlags()
	} else {
		cfg.ParseEnv()
	}
	cfg.CheckArgs()
	cfg.PrintArgs()
}
