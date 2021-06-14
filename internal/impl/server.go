package impl

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Telegram struct {
		BotToken    string `yaml:"bot_token"`
		AdminChatID int64  `yaml:"admin_chat_id"`
	} `yaml:"telegram"`
}

type Server struct {
	tg *telegramBotAPI
}

func NewServer() *Server {
	return &Server{}
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
	conf := getConf()
	crawler := NewCrawler()
	crawler.server = server
	server.tg, err = NewTelegramBot(conf.Telegram.AdminChatID, conf.Telegram.BotToken)
	if err != nil {
		return err
	}

	return crawler.Run()
}
