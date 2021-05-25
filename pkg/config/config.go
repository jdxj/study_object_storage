package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// PrintConfigFormat 用于生成格式化的配置文件.
// 如果有新配置被添加, 可以使用该函数自动生成 yaml
// 文件来防止手写 yaml 带来的错误.
func PrintConfigFormat() ([]byte, error) {
	data, _ := yaml.Marshal(Config{})

	fileName := "conf.yaml.default"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer file.Sync()

	_, err = file.Write(data)
	return data, err
}

func New(path string) (*Config, error) {
	conf := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return conf, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return conf, decoder.Decode(conf)
}

type Config struct {
	Logger Logger `yaml:"logger"`
	Gate   Gate   `yaml:"gate"`
	Web    Web    `yaml:"web"`
	MySQL  MySQL  `yaml:"mysql"`
	Mongo  Mongo  `yaml:"mongo"`
	Redis  Redis  `yaml:"redis"`
	Rabbit Rabbit `yaml:"rabbit"`
}

type Logger struct {
	FileName   string `yaml:"file_name"`
	AppName    string `yaml:"app_name"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Level      int    `yaml:"level"`
	LocalTime  bool   `yaml:"local_time"`
	Compress   bool   `yaml:"compress"`
}

type Gate struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Node int    `yaml:"node"`
}

type Web struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type MySQL struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Dbname string `yaml:"dbname"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type Mongo struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Dbname string `yaml:"dbname"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
}

type Redis struct {
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"db"`
}

type Rabbit struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
