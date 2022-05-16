package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MyConf struct {
	port        string
	db_url      string
	some_app_id string
	some_url    string
	checked     bool
	err         []error
	source      string
	path        string
}

func (cfg *MyConf) FlagParse() {
	args := flag.NewFlagSet("args", flag.ExitOnError)
	file := flag.NewFlagSet("file", flag.ExitOnError)
	args.StringVar(&cfg.port, "port", "", "type port number")
	args.StringVar(&cfg.db_url, "db_url", "", "type db_url")
	args.StringVar(&cfg.some_app_id, "some_app_id", "", "type some_app_id")
	args.StringVar(&cfg.some_url, "some_url", "", "type some url")
	file.StringVar(&cfg.path, "path", "", "type path to config file")
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "args":
			args.Parse(os.Args[2:])
			cfg.source = "args"
		case "file":
			file.Parse(os.Args[2:])
			cfg.LoadCfgFromFile()
		default:
			{
				fmt.Println("передайте первым параметром args или file, либо запустите утилиту без параметров, для чтения переменных окружения")
				os.Exit(1)
			}
		}

	} else {
		cfg.ParseEnv()
	}
	cfg.CheckArgs()
	cfg.PrintCfg()
}

func (cfg *MyConf) LoadCfgFromFile() {
	cfgMap := make(map[string]string)
	file, err := ioutil.ReadFile(cfg.path)
	if err != nil {
		fmt.Println("Укажите верный путь для файла конфигурации")
		os.Exit(2)
	}
	lines := strings.Split(string(file), "\n")
	for i := range lines {
		temp := strings.SplitN(lines[i], ":", 2)
		if len(temp) == 2 {
			cfgMap[temp[0]] = strings.ReplaceAll((temp[1]), " ", "")
		}
	}
	cfg.port = cfgMap["port"]
	cfg.db_url = cfgMap["db_url"]
	cfg.some_url = cfgMap["some_url"]
	cfg.some_app_id = cfgMap["some_app_id"]
	cfg.source = "file"

}

func (cfg *MyConf) ParseEnv() {
	cfg.port = os.Getenv("port")
	cfg.db_url = os.Getenv("db_url")
	cfg.some_app_id = os.Getenv("some_app_id")
	cfg.some_url = os.Getenv("some_url")
	cfg.source = "env"
}

func (cfg *MyConf) PrintCfg() {
	if cfg.checked {
		fmt.Printf("%s\n%s\n%s\n%s\n", cfg.port, cfg.db_url, cfg.some_app_id, cfg.some_url)
		switch cfg.source {
		case "env":
			{
				fmt.Println("Источник данных: переменные окружения")
			}
		case "args":
			{
				fmt.Println("Источник данных: аргументы командной строки")
			}
		case "file":
			{
				fmt.Println("Источник данных: файл конфигурации")
			}
		}

	} else {
		for i := 0; i < len(cfg.err); i++ {
			fmt.Println(cfg.err[i])
		}
	}
}

func (cfg *MyConf) CheckArgs() {
	if cfg.some_app_id == "" {
		cfg.err = append(cfg.err, errors.New("не указан идентификатор приложения [a-Z, 0-9]"))
	}
	_, portErr := strconv.Atoi(cfg.port)
	if portErr != nil {
		cfg.err = append(cfg.err, errors.New("не указан номер порта [0-9]{1-6}"))
	}
	_, ok := url.ParseRequestURI(cfg.some_url)
	if ok != nil {
		cfg.err = append(cfg.err, errors.New("не указан валидный url"))
	}
	if len(cfg.err) == 0 {
		cfg.checked = true
	}
	re, _ := regexp.Compile(`^postgres://.+:.+@.+:\d{0,5}/.+\?sslmode=(enable$|disable$)`)
	if !re.MatchString(cfg.db_url) {
		cfg.err = append(cfg.err, errors.New("не указана валидная строка подключения к БД, например [postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=disable]"))
	}
	if len(cfg.err) == 0 {
		cfg.checked = true
	} else {
		cfg.checked = false
	}
}
