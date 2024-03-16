package config

import (
	"gotemplate/logger"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance Econfig
)

type config struct {
	AppName            string
	AppEnv             string
	DBConnection       string
	TokenSymmetricKey  string
	HttpUrl            string
	HttpPort           string
	DBHost             string
	DBPort             string
	DBdatabase         string
	DBUsername         string
	DBPassword         string
	TokenDuration      string
	RedisServer        string
	RedisPassword      string
	HttpAllowedOrigins string
	Loglevel           string
	ShutDownTime       string
	ShutDowntype       string

	MaxConns        int
	MinConns        int
	MaxConnLifetime int
	MaxConnIdleTime int
}

/*
type Econfig struct {
	appName            string `yaml:"AppName"`
	appEnv             string `yaml:"AppEnv"`
	dBConnection       string `yaml:"DBConnection"`
	tokenSymmetricKey  string `yaml:"TokenSymmetricKey"`
	httpUrl            string `yaml:"HttpUrl"`
	httpPort           string `yaml:"HttpPort"`
	dBHost             string `yaml:"DBHost"`
	dBPort             string `yaml:"DBPort"`
	dBdatabase         string `yaml:"DBdatabase"`
	dBUsername         string `yaml:"DBUsername"`
	dBPassword         string `yaml:"DBPassword"`
	tokenDuration      string `yaml:"TokenDuration"`
	redisServer        string `yaml:"RedisServer"`
	redisPassword      string `yaml:"RedisPassword"`
	httpAllowedOrigins string `yaml:"HttpAllowedOrigins"`
	loglevel           string `yaml:"Loglevel"`
	shutDownTime       string `yaml:"ShutDownTime"`
	shutDowntype       string `yaml:"ShutDowntype"`
}*/

type Econfig struct {
	appName            string `mapstructure:"AppName"`
	appEnv             string `mapstructure:"AppEnv"`
	dBConnection       string `mapstructure:"DBConnection"`
	tokenSymmetricKey  string `mapstructure:"TokenSymmetricKey"`
	httpUrl            string `mapstructure:"HttpUrl"`
	httpPort           string `mapstructure:"HttpPort"`
	dBHost             string `mapstructure:"DBHost"`
	dBPort             string `mapstructure:"DBPort"`
	dBdatabase         string `mapstructure:"DBdatabase"`
	dBUsername         string `mapstructure:"DBUsername"`
	dBPassword         string `mapstructure:"DBPassword"`
	tokenDuration      string `mapstructure:"TokenDuration"`
	redisServer        string `mapstructure:"RedisServer"`
	redisPassword      string `mapstructure:"RedisPassword"`
	httpAllowedOrigins string `mapstructure:"HttpAllowedOrigins"`
	loglevel           string `mapstructure:"Loglevel"`
	shutDownTime       string `mapstructure:"ShutDownTime"`
	shutDowntype       string `mapstructure:"ShutDowntype"`
	maxConns           int    `mapstructure:"MaxConns"`
	minConns           int    `mapstructure:"MinConns"`
	maxConnLifetime    int    `mapstructure:"MaxConnLifetime"`
	maxConnIdleTime    int    `mapstructure:"MaxConnIdleTime"`
}

func NewConfig(c config) Econfig {
	return Econfig{
		appName:            c.AppName,
		appEnv:             c.AppEnv,
		dBConnection:       c.DBConnection,
		tokenSymmetricKey:  c.TokenSymmetricKey,
		httpUrl:            c.HttpUrl,
		httpPort:           c.HttpPort,
		dBHost:             c.DBHost,
		dBPort:             c.DBPort,
		dBdatabase:         c.DBdatabase,
		dBUsername:         c.DBUsername,
		dBPassword:         c.DBPassword,
		tokenDuration:      c.TokenDuration,
		redisServer:        c.RedisServer,
		redisPassword:      c.RedisPassword,
		httpAllowedOrigins: c.HttpAllowedOrigins,
		loglevel:           c.Loglevel,
		shutDownTime:       c.ShutDownTime,
		shutDowntype:       c.ShutDowntype,
		maxConns:           c.MaxConns,
		minConns:           c.MinConns,
		maxConnLifetime:    c.MaxConnLifetime,
		maxConnIdleTime:    c.MaxConnIdleTime,
	}
}

func Load(log *logger.Logger) Econfig {
	c := config{}
	once.Do(func() {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")

		//log := logger.New("DEBUG")

		// Attempt to read the configuration file
		if err := viper.ReadInConfig(); err != nil {
			log.Error("Failed to read configuration", err)
			os.Exit(1)
		}

		// Unmarshal the configuration into the struct
		if err := viper.Unmarshal(&c); err != nil {
			log.Error("Failed to unmarshal configuration", err)
			os.Exit(1)
		}

		instance = NewConfig(c)

		// Assign configuration values to environment variables
		// os.Setenv("APP_NAME", c.AppName)
		// os.Setenv("APP_ENV", c.AppEnv)
		// os.Setenv("DB_CONNECTION", c.DBConnection)
		// os.Setenv("TOKEN_SYMMETRIC_KEY", c.TokenSymmetricKey)
		// os.Setenv("HTTP_URL", c.HttpUrl)
		// os.Setenv("HTTP_PORT", c.HttpPort)
		// os.Setenv("DB_HOST", c.DBHost)
		// os.Setenv("DB_PORT", c.DBPort)
		// os.Setenv("DB_DATABASE", c.DBdatabase)
		// os.Setenv("DB_USERNAME", c.DBUsername)
		// os.Setenv("DB_PASSWORD", c.DBPassword)
		// os.Setenv("HTTP_ALLOWED_ORIGINS", c.HttpAllowedOrigins)
		// os.Setenv("REDIS_SERVER", c.RedisServer)
		// os.Setenv("REDIS_PASSWORD", c.RedisPassword)
		// os.Setenv("TOKEN_DURATION", c.TokenDuration)
		// os.Setenv("LOG_LEVEL", c.Loglevel)
		// os.Setenv("SHUTDOWN_TIME", c.ShutDownTime)
		// os.Setenv("SHUTDOWN_TYPE", c.ShutDowntype)
	})
	return instance
}

func (c *Econfig) AppName() string {
	return c.appName

}

func (c *Econfig) AppEnv() string {
	return c.appEnv

}
func (c *Econfig) DBConnection() string {
	return c.dBConnection

}

func (c *Econfig) TokenSymmetricKey() string {
	return c.tokenSymmetricKey

}

func (c *Econfig) Dbhost() string {
	return c.dBHost

}

// HttpUrl returns the httpUrl field value.
func (c *Econfig) HttpUrl() string {
	return c.httpUrl
}

// HttpPort returns the httpPort field value.
func (c *Econfig) HttpPort() string {
	return c.httpPort
}

// DBHost returns the dBHost field value.
func (c *Econfig) DBHost() string {
	return c.dBHost
}

// DBPort returns the dBPort field value.
func (c *Econfig) DBPort() string {
	return c.dBPort
}

// DBDatabase returns the dBdatabase field value.
func (c *Econfig) DBDatabase() string {
	return c.dBdatabase
}

// DBUsername returns the dBUsername field value.
func (c *Econfig) DBUsername() string {
	return c.dBUsername
}

// DBPassword returns the dBPassword field value.
func (c *Econfig) DBPassword() string {
	return c.dBPassword
}

// TokenDuration returns the tokenDuration field value.
func (c *Econfig) TokenDuration() string {
	return c.tokenDuration
}

// RedisServer returns the redisServer field value.
func (c *Econfig) RedisServer() string {
	return c.redisServer
}

// RedisPassword returns the redisPassword field value.
func (c *Econfig) RedisPassword() string {
	return c.redisPassword
}

// HttpAllowedOrigins returns the httpAllowedOrigins field value.
func (c *Econfig) HttpAllowedOrigins() string {
	return c.httpAllowedOrigins
}

// LogLevel returns the loglevel field value.
func (c *Econfig) LogLevel() string {
	return c.loglevel
}

// ShutDownTime returns the shutDownTime field value.
func (c *Econfig) ShutDownTime() string {
	return c.shutDownTime
}

// ShutDownType returns the shutDowntype field value.
func (c *Econfig) ShutDownType() string {
	return c.shutDowntype
}

func (c *Econfig) MaxConns() int {
	return c.maxConns
}

func (c *Econfig) MinConns() int {
	return c.minConns
}

func (c *Econfig) MaxConnLifetime() int {
	return c.maxConnLifetime
}

func (c *Econfig) MaxConnIdleTime() int {
	return c.maxConnIdleTime
}
