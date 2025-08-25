package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type ApiConfig struct {
	Prefix string `yaml:"prefix"`
}

// JWT 配置结构体
type JWTConfig struct {
	SecretKey      string `yaml:"secretKey"`      // JWT 的密钥
	ExpirationTime string `yaml:"expirationTime"` // JWT 的过期时间
	Issuer         string `yaml:"issuer"`         // JWT 的发行者
	Audience       string `yaml:"audience"`       // JWT 的受众
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
	DialTimeout  string `yaml:"dialTimeout"`
	ReadTimeout  string `yaml:"readTimeout"`
	WriteTimeout string `yaml:"writeTimeout"`
}

type RateConfig struct {
	UserLimit int `yaml:"userLimit"`
	ApiLimit  int `yaml:"apiLimit"`
}

type RabbitmqConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MqConfig struct {
	Exchange   string `yaml:"exchange"`
	Queue      string `yaml:"queue"`
	RoutingKey string `yaml:"routingKey"`
	Handler    string `yaml:"handler"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	BaseUrl   string `yaml:"baseUrl"`
}
type WebSocketConfig struct {
	Addr string `yaml:"addr"`
}

// Config 配置结构体 整个文件
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Api       ApiConfig       `yaml:"api"`
	Jwt       JWTConfig       `yaml:"jwt"`
	Redis     RedisConfig     `yaml:"redis"`
	Rate      RateConfig      `yaml:"rate"`
	Rabbitmq  RabbitmqConfig  `yaml:"rabbitmq"`
	Mq        []MqConfig      `yaml:"mq"`
	Minio     MinioConfig     `yaml:"minio"`
	WebSocket WebSocketConfig `yaml:"websocket"`
}

var AppConfig *Config

// LoadConfig 加载配置文件，根据环境选择加载不同的配置文件
func LoadConfig() error {
	// 重置配置
	AppConfig = &Config{}
	configPath := os.Getenv("CONFIG_PATH") // 通过环境变量获取配置路径
	if configPath == "" {
		configPath = "./configs" // 默认路径
	}
	// 设置默认配置文件路径和文件名
	viper.SetConfigName("app")      // 默认的配置文件
	viper.AddConfigPath(configPath) // 配置文件所在路径
	viper.SetConfigType("yaml")     // 配置文件类型

	// 加载默认配置文件 app.yaml
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error loading default config file: %v", err)
	}
	fmt.Println("Loading default config from:", viper.ConfigFileUsed())

	// 获取环境变量，并加载对应的环境配置文件（如 app.test.yaml）
	env := os.Getenv("APP_ENV")
	if env != "" {
		// 根据环境变量加载特定的配置文件
		viper.SetConfigName(fmt.Sprintf("app.%s", env)) // app.test.yaml, app.dev.yaml
		// 再次添加配置路径，覆盖同名字段
		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("error loading environment config file: %v", err)
		}
		fmt.Printf("Loaded config for environment: %s\n", env)
	} else {
		fmt.Println("No APP_ENV variable set, using default config.")
	}

	// 配置文件变化监听
	viper.WatchConfig()

	// 配置文件变化时的回调函数
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err := LoadConfig()
		if err != nil {
			log.Printf("Error reloading config: %v", err)
		}
	})

	// Unmarshal 配置到结构体
	if err := viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("unable to unmarshal config: %v", err)
	}

	// 输出加载的配置
	fmt.Printf("Loading config from: %s\n", viper.ConfigFileUsed())
	return nil
}
