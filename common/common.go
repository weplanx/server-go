package common

import (
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/encryption"
	"github.com/weplanx/go/passport"
	"go.mongodb.org/mongo-driver/mongo"
)

const TokenClaimsKey = "token-claims"

type Inject struct {
	Values      *Values
	MongoClient *mongo.Client
	Db          *mongo.Database
	Redis       *redis.Client
	Nats        *nats.Conn
	Js          nats.JetStreamContext
	Passport    *passport.Passport
	Cipher      *encryption.Cipher
	Idx         *encryption.IDx
}

type Values struct {
	Address        string          `yaml:"address"`
	TrustedProxies []string        `yaml:"trusted_proxies"`
	Name           string          `yaml:"name"`
	Key            string          `yaml:"key"`
	Cors           Cors            `yaml:"cors"`
	Database       Database        `yaml:"database"`
	Redis          Redis           `yaml:"redis"`
	Nats           Nats            `yaml:"nats"`
	Passport       passport.Option `yaml:"passport"`
}

type Cors struct {
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type Database struct {
	Uri    string `yaml:"uri"`
	DbName string `yaml:"dbName"`
}

type Redis struct {
	Uri string `yaml:"uri"`
}

type Nats struct {
	Hosts []string `yaml:"hosts"`
	Nkey  string   `yaml:"nkey"`
}

func (x *Values) RedisKey(name string) string {
	return x.Name + ":" + name
}
