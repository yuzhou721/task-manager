package conf

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

//Database 数据库相关配置
type Database struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	DbName   string `yaml:"dbName"`
}

type Server struct {
	HTTPPort     int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type Log struct {
	Path   string `yaml:"path"`
	Prefix string `yaml:"prefix"`
}

type Setting struct {
	RunMode  string   `yaml:"runMode"`
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Log      Log      `yaml:"log"`
}

//Config 全局配置文件
var Config = &Setting{}

func init() {
	path := "etc/" + os.Getenv("ENV") + ".yaml"
	log.Println(path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error load config file")
	}

	err = yaml.Unmarshal(yamlFile, Config)
	if err != nil {
		log.Fatalln("error load yaml file")
	}
	Config.Server.ReadTimeout = Config.Server.ReadTimeout * time.Second
	Config.Server.WriteTimeout = Config.Server.WriteTimeout * time.Second
}
