package impl

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-redis/redis/v8"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Telegram struct {
		BotToken    string `yaml:"bot_token"`
		AdminChatID int64  `yaml:"admin_chat_id"`
	} `yaml:"telegram"`
	Redis struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"redis"`
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
	server.tg, err = NewTelegramBot(server.config.Telegram.AdminChatID, server.config.Telegram.BotToken)
	if err != nil {
		return err
	}

	server.db = NewRedisConnect(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", server.config.Redis.Host, server.config.Redis.Port),
	})

	crawler := NewCrawler(server)

	return crawler.Run(context.Background())
}
