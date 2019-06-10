package bootstarp

import (
	"github.com/spf13/viper"
	"os"
	"fmt"
)

var Config *viper.Viper
var configFile string

func SetConfigFile(filename string) {
	configFile = filename
}

func InitConfig() (*viper.Viper  ,error){
	Config = viper.New()
	if configFile != "" {
		_, err := os.Stat(configFile)
		if err == nil {
			Config.SetConfigFile(configFile)
		} else {
			fmt.Errorf("配置文件不存在:%s", configFile)
		}
	} else {
		Config.SetConfigType("yaml")            // or Config.SetConfigType("YAML")
		Config.SetConfigName("config")          // name of config file (without extension)
		Config.AddConfigPath("/etc/GMQ/")  // path to look for the config file in
		Config.AddConfigPath("$HOME/.GMQ") // call multiple times to add many search paths
		Config.AddConfigPath(".")
	}

	err := Config.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		return nil ,err
	}
	return Config ,nil
}

func GetConfig() *viper.Viper {
	if Config != nil {
		return Config
	}
	InitConfig()
	return Config
}
