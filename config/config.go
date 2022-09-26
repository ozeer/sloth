package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"time"
)

type Conf struct {
	Gin           SectionGin           `ini:"gin"`
	Core          SectionCore          `ini:"core"`
	Log           SectionLog           `ini:"log"`
	Redis         SectionRedis         `ini:"redis"`
	DingTalk      SectionDingTalk      `int:"ding_talk"`
	ConsumerQueue SectionConsumerQueue `ini:"consumer_queue"`
}

type SectionGin struct {
	Mode    string `ini:"mode"`
	Address string `ini:"address"`
	Port    string `ini:"port"`
}

type SectionCore struct {
	BucketSize        int    `ini:"bucket_size"`
	BucketName        string `ini:"bucket_name"`
	QueueName         string `ini:"queue_name"`
	QueueBlockTimeout int64  `ini:"queue_block_timeout"`
}

type SectionLog struct {
	Format      string `ini:"format"`
	AccessLog   string `ini:"access_log"`
	AccessLevel string `ini:"access_level"`
	ErrorLog    string `ini:"error_log"`
	ErrorLevel  string `ini:"error_level"`
}

type SectionRedis struct {
	Host         string        `ini:"host"`
	Port         string        `ini:"port"`
	Password     string        `ini:"password"`
	DB           int           `ini:"db"`
	DialTimeout  time.Duration `ini:"dial_timeout"`
	ReadTimeout  time.Duration `ini:"read_timeout"`
	WriteTimeout time.Duration `ini:"write_timeout"`
}

type SectionDingTalk struct {
	SwitchOn    bool   `ini:"switch_on"`
	BaseUrl     string `ini:"base_url"`
	AccessToken string `ini:"access_token"`
}

type SectionConsumerQueue struct {
	Host      string `ini:"host"`
	Port      string `ini:"port"`
	Password  string `ini:"password"`
	QueueName string `ini:"queue_name"`
}

var defaultConf = []byte(`
[gin]
  mode=debug
  address=0.0.0.0
  port=9288
[core]
  bucket_size=3
  bucket_name=dq_bucket_%d
  queue_name=dq_queue_%s
  queue_block_timeout=178
[log]
  format="string"
  access_log="stdout"
  access_level="debug"
  error_log="stderr"
  error_level="error"
[redis]
  host="127.0.0.1"
  port="6379"
  password=""
  db=0
  connect_timeout=5000
  read_timeout=180000
  write_timeout=3000
[ding_talk]
  base_url=https://oapi.dingtalk.com/robot/send?access_token=
  access_token=""
  switch_on=false
[consumer_queue]
	host="127.0.0.1"
	port="6379"
	password=""
	queue_name="hzQDoraemon"
`)

func LoadConfig(confPath string) (Conf, error) {
	var conf Conf

	viper.SetConfigType("ini")
	viper.AutomaticEnv()

	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)
		if err != nil {
			return conf, err
		}
		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")

		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Println("Using default config file")
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	// gin
	conf.Gin.Mode = viper.GetString("gin.mode")
	conf.Gin.Address = viper.GetString("gin.address")
	conf.Gin.Port = viper.GetString("gin.port")

	// core
	conf.Core.BucketSize = viper.GetInt("core.bucket_size")
	conf.Core.BucketName = viper.GetString("core.bucket_name")
	conf.Core.QueueName = viper.GetString("core.queue_name")
	conf.Core.QueueBlockTimeout = viper.GetInt64("core.queue_block_timeout")

	// log
	conf.Log.Format = viper.GetString("log.format")
	conf.Log.AccessLog = viper.GetString("log.access_log")
	conf.Log.AccessLevel = viper.GetString("log.access_level")
	conf.Log.ErrorLog = viper.GetString("log.error_log")
	conf.Log.ErrorLevel = viper.GetString("log.error_level")

	// redis
	conf.Redis.Host = viper.GetString("redis.host")
	conf.Redis.Port = viper.GetString("redis.port")
	conf.Redis.Password = viper.GetString("redis.password")
	conf.Redis.DB = viper.GetInt("redis.db")

	// consumer queue
	conf.ConsumerQueue.Host = viper.GetString("consumer_queue.host")
	conf.ConsumerQueue.Port = viper.GetString("consumer_queue.port")
	conf.ConsumerQueue.Password = viper.GetString("consumer_queue.password")
	conf.ConsumerQueue.QueueName = viper.GetString("consumer_queue.queue_name")

	// ding talk
	conf.DingTalk.BaseUrl = viper.GetString("ding_talk.base_url")
	conf.DingTalk.AccessToken = viper.GetString("ding_talk.access_token")
	conf.DingTalk.SwitchOn = viper.GetBool("ding_talk.switch_on")

	return conf, nil
}
