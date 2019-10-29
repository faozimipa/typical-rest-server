package typredis

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"

	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/urfave/cli"
	"go.uber.org/dig"
)

// Module of redis
func Module() interface{} {
	return &redisModule{
		Name: "Redis",
		Configuration: typictx.Configuration{
			Prefix: "REDIS",
			Spec:   &Config{},
		},
	}
}

type redisModule struct {
	typictx.Configuration
	Name string
}

// Construct dependencies
func (r redisModule) Construct(c *dig.Container) (err error) {
	return c.Provide(r.connect)
}

// Destruct dependencies
func (r redisModule) Destruct(c *dig.Container) (err error) {
	return c.Invoke(c)
}

// CommandLine return command
func (r redisModule) CommandLine() cli.Command {
	return cli.Command{
		Name:   "redis",
		Usage:  "Redis Utility Tool",
		Before: r.cliBefore,
		Subcommands: []cli.Command{
			{Name: "console", ShortName: "c", Action: r.action(r.console)},
		},
	}
}

func (redisModule) cliBefore(ctx *cli.Context) (err error) {
	return typienv.LoadEnvFile()
}

func (r redisModule) action(fn interface{}) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		c := dig.New()
		if err = c.Provide(r.loadConfig); err != nil {
			return
		}
		return c.Invoke(fn)
	}
}

func (r redisModule) loadConfig() (cfg *Config, err error) {
	cfg = new(Config)
	err = envconfig.Process(r.Configure().Prefix, cfg)
	return
}

func (redisModule) connect(cfg *Config) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:           cfg.Password,
		DB:                 cfg.DB,
		PoolSize:           cfg.PoolSize,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadWriteTimeout,
		WriteTimeout:       cfg.ReadWriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		MaxConnAge:         cfg.MaxConnAge,
	})
	err = client.Ping().Err()
	return
}

func (redisModule) disconnect(client *redis.Client) (err error) {
	return client.Close()
}

func (redisModule) console(config *Config) (err error) {
	args := []string{
		"-h", config.Host,
		"-p", config.Port,
	}
	if config.Password != "" {
		args = append(args, "-a", config.Password)
	}
	// TODO: using docker -it
	cmd := exec.Command("redis-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}