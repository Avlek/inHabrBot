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
	URL     string `yaml:"url"`
	Timeout int64  `yaml:"timeout"`
}

type Config struct {
	Telegram struct {
		BotToken  string `yaml:"bot_token"`
		ChannelID int64  `yaml:"channel_id"`
	} `yaml:"telegram"`
	Redis  RedisConfig  `yaml:"redis"`
	Parser ParserConfig `yaml:"parser"`
}

type Server struct {
	config *Config
	tg     *telegramBotAPI
	db     *DB
}

func NewServer() *Server {
	return &Server{
		config: getConf(),
	}
}

func getConf() *Config {
	yamlFile, err := ioutil.ReadFile("configs/dev.yaml")
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

func (server *Server) Run() (err error) {
	server.db = NewRedisConnect(server.config.Redis)

	tg, err := NewTelegramBot(server.config.Telegram.ChannelID, server.config.Telegram.BotToken)
	if err != nil {
		return err
	}
	server.tg = tg

	ctx := context.Background()

	crawler := NewCrawler(server)
	err = crawler.InitCrawler(ctx)
	if err != nil {
		return err
	}

	return crawler.Run(ctx)
}
