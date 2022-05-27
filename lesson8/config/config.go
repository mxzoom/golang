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

	"gopkg.in/yaml.v2"
)

type MyConf struct {
	Port        string `yaml:"port"`
	Db_url      string `yaml:"db_url"`
	Some_app_id string `yaml:"some_app_id"`
	Some_url    string `yaml:"some_url"`
	checked     bool
	err         []error
	source      string
	path        string
}

func (cfg *MyConf) FlagParse() {
	args := flag.NewFlagSet("args", flag.ExitOnError)
	file := flag.NewFlagSet("file", flag.ExitOnError)
	args.StringVar(&cfg.Port, "port", "", "type port number")
	args.StringVar(&cfg.Db_url, "db_url", "", "type db_url")
	args.StringVar(&cfg.Some_app_id, "some_app_id", "", "type some_app_id")
	args.StringVar(&cfg.Some_url, "some_url", "", "type some url")
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
	file, err := ioutil.ReadFile(cfg.path)
	if err != nil {
		fmt.Println("Укажите верный путь для файла конфигурации")
		os.Exit(2)
	}
	yamlErr := yaml.Unmarshal(file, cfg)
	if yamlErr != nil {
		fmt.Println("Файл конфигурации отсутсвует, или имеет неизвестный формат")
		os.Exit(2)
	}

}

func (cfg *MyConf) ParseEnv() {
	cfg.Port = os.Getenv("port")
	cfg.Db_url = os.Getenv("db_url")
	cfg.Some_app_id = os.Getenv("some_app_id")
	cfg.Some_url = os.Getenv("some_url")
	cfg.source = "env"
}

func (cfg *MyConf) PrintCfg() {
	if cfg.checked {
		fmt.Printf("%s\n%s\n%s\n%s\n", cfg.Port, cfg.Db_url, cfg.Some_app_id, cfg.Some_url)
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
	if cfg.Some_app_id == "" {
		cfg.err = append(cfg.err, errors.New("не указан идентификатор приложения [a-Z, 0-9]"))
	}
	_, portErr := strconv.Atoi(cfg.Port)
	if portErr != nil {
		cfg.err = append(cfg.err, errors.New("не указан номер порта [0-9]{1-6}"))
	}
	_, ok := url.ParseRequestURI(cfg.Some_url)
	if ok != nil {
		cfg.err = append(cfg.err, errors.New("не указан валидный url"))
	}
	if len(cfg.err) == 0 {
		cfg.checked = true
	}
	re, _ := regexp.Compile(`^postgres://.+:.+@.+:\d{0,5}/.+\?sslmode=(enable$|disable$)`)
	if !re.MatchString(cfg.Db_url) {
		cfg.err = append(cfg.err, errors.New("не указана валидная строка подключения к БД, например [postgres://db-user:db-password@petstore-db:5432/petstore?sslmode=disable]"))
	}
	if len(cfg.err) == 0 {
		cfg.checked = true
	} else {
		cfg.checked = false
	}
}
