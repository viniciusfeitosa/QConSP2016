package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/rakyll/globalconf"
)

// Config is the struct with configs to the app
type Config struct {
	Address    string
	Debug      bool
	Mongo      MongoConfig
	Redis      RedisConfig
	NumWorkers NumWorkersConfig
}

// MongoConfig is the struct with configs of mongo
type MongoConfig struct {
	Address1     string
	DatabaseName string
}

// RedisConfig is the struct with configs of redis
type RedisConfig struct {
	Address         string
	Auth            string
	DB              uint
	MaxIdle         uint
	MaxActive       uint
	IdleTimeoutSecs uint
}

// NumWorkersConfig is the struct with information of numbers of workers run.
type NumWorkersConfig struct {
	NumUsersMongoWorkers int
}

var fs = &flag.FlagSet{}

// Print error, usage and exit with code
func printErrorUsageAndExitWithCode(err string, code int) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
	printUsage()
	os.Exit(code)
}

// Print command line help
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [flags] [CONFIG]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fs.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	fmt.Fprintf(os.Stderr, "  CONFIG: Config file path\n")
	fmt.Fprintf(os.Stderr, "\n")
}

// LoadConfig parse command line and load config file
func LoadConfig() *Config {
	cfg := &Config{}

	fs.StringVar(&cfg.Address, "address", ":8888", "Address where run the app")
	fs.BoolVar(&cfg.Debug, "debug", false, "Enable debugging")
	fs.StringVar(&cfg.Mongo.Address1, "mongo_address1", "localhost", "host to connect on Mongo")
	fs.StringVar(&cfg.Mongo.DatabaseName, "mongo_databasename", "qconsp", "Database name for connect on Mongo")
	fs.StringVar(&cfg.Redis.Address, "redis_address", "localhost:6379", "Redis server address")
	// fs.StringVar(&cfg.Redis.Auth, "redis_auth", "", "Redis authentication token")
	fs.UintVar(&cfg.Redis.DB, "redis_db", 0, "Redis database number")
	fs.UintVar(&cfg.Redis.MaxIdle, "redis_max_idle", 10, "Max idle conns to Redis")
	fs.UintVar(&cfg.Redis.MaxActive, "redis_max_active", 1000, "Max active conns to Redis")
	fs.UintVar(&cfg.Redis.IdleTimeoutSecs, "redis_idle_timeout_secs", 60, "Redis idle conn timeout in seconds")
	fs.IntVar(&cfg.NumWorkers.NumUsersMongoWorkers, "num_users_mongo_workers", 10, "Number of the workers to UsersMongoWorkers")

	fs.Usage = printUsage
	fs.Parse(os.Args[1:])
	if fs.NArg() != 1 {
		printErrorUsageAndExitWithCode("Wrong or missing arguments", 1)
	}

	gcfg, err := globalconf.NewWithOptions(&globalconf.Options{Filename: fs.Arg(0), EnvPrefix: "QCONSP_"})
	if err != nil {
		printErrorUsageAndExitWithCode(err.Error(), 2)

	}

	globalconf.Register("", fs)
	gcfg.ParseAll()

	return cfg
}
