package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"time"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}

	parseInt := func(key string) int {
		value, err := strconv.Atoi(envMap[key])
		if err != nil {
			log.Fatalf("load %s failed: %v", key, err)
		}
		return value
	}

	parseDuration := func(key string) time.Duration {
		value := parseInt(key)
		return time.Duration(value) * time.Second
	}

	return &config{
		app: &app{
			host:         envMap["APP_HOST"],
			port:         parseInt("APP_PORT"),
			name:         envMap["APP_NAME"],
			version:      envMap["APP_VERSION"],
			readTimeout:  parseDuration("APP_READ_TIMEOUT"),
			writeTimeout: parseDuration("APP_WRITE_TIMEOUT"),
		},
		db: &db{
			host:           envMap["DB_HOST"],
			port:           parseInt("DB_PORT"),
			protocol:       envMap["DB_PROTOCOL"],
			username:       envMap["DB_USERNAME"],
			password:       envMap["DB_PASSWORD"],
			database:       envMap["DB_DATABASE"],
			maxConnections: parseInt("DB_MAX_CONNECTIONS"),
			mongoURI:       envMap["MONGO_URI"],
		},
	}
}

type IConfig interface {
	App() IAppConfig
	Db() IDbConfig
}

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
}

func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }

type IDbConfig interface {
	Url() string
	MaxOpenConns() int
	MongoURI() string
}

func (d *db) Url() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", d.username, d.password, d.host, d.port, d.database)
}

func (d *db) MongoURI() string {
	return d.mongoURI
}

func (d *db) MaxOpenConns() int {
	return d.maxConnections
}

type config struct {
	app *app
	db  *db
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	maxConnections int
	mongoURI       string
}

func (c *config) App() IAppConfig {
	return c.app
}

func (c *config) Db() IDbConfig {
	return c.db
}
