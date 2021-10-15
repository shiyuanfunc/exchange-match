package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
	NameSrv              []string `json:"nameSrv"`
	TradeConsumerGroup   string
	DefaultConsumerGroup string
	ProducerGroup        string
}

var configFile ConfigFile

const (
	TEST_TOPIC       = "shiyuan-test"
	PROFILE          = "profile"
	CONFIG_FILE_NAME = "config%s.json"
)

func init() {
	configFileName := fmt.Sprintf(CONFIG_FILE_NAME, "")
	profile := os.Getenv(PROFILE)
	if profile != "" {
		configFileName = fmt.Sprintf(CONFIG_FILE_NAME, "-"+profile)
	}
	file, err := ioutil.ReadFile(configFileName)
	if err != nil {
		os.Exit(-1)
	}
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		os.Exit(-1)
	}
}

// 读取配置文件
func main() {
	config := GetConfig()
	fmt.Println("config content >>>> ", (*config).TradeConsumerGroup)
}

func GetConfig() *ConfigFile {
	return &configFile
}

// rocketmq 消费组
func GetConsumerGroup() string {
	config := GetConfig()
	return (*config).TradeConsumerGroup
}

// rocketmq 的nameSrv
func GetRocketNameSrv() []string {
	config := GetConfig()
	return (*config).NameSrv
}
