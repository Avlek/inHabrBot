package impl

import (
	"context"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       int64  `yaml:"port"`
	Expiration int64  `yaml:"expiration"`
}

type ParserConfig struct {
	URLS    []string `yaml:"urls"`
	Version int8     `yaml:"version"`
	Timeout int64    `yaml:"timeout"`
}

type TelegramConfig struct {
	BotToken  string `yaml:"bot_token"`
	ChannelID int64  `yaml:"channel_id"`
	AdminID   int64  `yaml:"admin_id"`
}

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Redis    RedisConfig    `yaml:"redis"`
	Parser   ParserConfig   `yaml:"parser"`
}

type Server struct {
	config *Config
	tg     *telegramBotAPI
	db     *DB
}

func NewServer(fileName string) *Server {
	return &Server{
		config: getConf(fileName),
	}
}

func getConf(fileName string) *Config {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	c := Config{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &c
}

func (server *Server) Run(toInit bool) (err error) {
	server.db = NewRedisConnect(server.config)

	tg, err := NewTelegramBot(server.config.Telegram)
	if err != nil {
		return err
	}
	server.tg = tg

	ctx := context.Background()

	crawler := NewCrawler(server)
	if toInit {
		err = crawler.InitCrawler(ctx)
		if err != nil {
			return err
		}
	}

	return crawler.Run(ctx)
}
