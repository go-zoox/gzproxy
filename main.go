package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzproxy/core"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "gzproxy",
		Usage:   "gzproxy is a portable proxy cli",
		Version: Version,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
			&cli.StringFlag{
				Name:     "upstream",
				Usage:    "upstream service",
				EnvVars:  []string{"UPSTREAM"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "prefix",
				Usage:    "the prefix",
				EnvVars:  []string{"PREFIX"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "basic-username",
				Usage:   "basic username",
				EnvVars: []string{"BASIC_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "basic-password",
				Usage:   "basic password",
				EnvVars: []string{"BASIC_PASSWORD"},
			},
			&cli.StringFlag{
				Name:    "auth-service",
				Usage:   "auth service",
				EnvVars: []string{"AUTH_SERVICE"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		basicUsername := ctx.String("basic-username")
		basicUpassword := ctx.String("basic-password")
		authService := ctx.String("auth-service")

		return core.Serve(&core.Config{
			Port:          ctx.Int64("port"),
			Upstream:      ctx.String("upstream"),
			Prefix:        ctx.String("prefix"),
			BasicUsername: basicUsername,
			BasicPassword: basicUpassword,
			AuthService:   authService,
		})
	})

	app.Run()
}
