package main

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/strings"
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
				Name:    "prefix",
				Usage:   "the prefix",
				EnvVars: []string{"PREFIX"},
			},
			&cli.BoolFlag{
				Name:    "disable-change-origin",
				Usage:   "disable change origin, default is false",
				EnvVars: []string{"DISABLE_CHANGE_ORIGIN"},
			},
			// basic auth
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
			// bearer token
			&cli.StringFlag{
				Name:    "bearer-token",
				Usage:   "bearer token",
				EnvVars: []string{"BEARER_TOKEN"},
			},
			// auth service
			&cli.StringFlag{
				Name:    "auth-service",
				Usage:   "auth service",
				EnvVars: []string{"AUTH_SERVICE"},
			},
			// oauth2
			&cli.StringFlag{
				Name:    "oauth2-provider",
				Usage:   "oauth2 provider, support: doreamon, github, feishu",
				EnvVars: []string{"OAUTH2_PROVIDER"},
			},
			&cli.StringFlag{
				Name:    "oauth2-client-id",
				Usage:   "oauth2 client id",
				EnvVars: []string{"OAUTH2_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "oauth2-client-secret",
				Usage:   "oauth2 client secret",
				EnvVars: []string{"OAUTH2_CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "oauth2-redirect-uri",
				Usage:   "oauth2 redirect uri",
				EnvVars: []string{"OAUTH2_REDIRECT_URI"},
			},
			//
			&cli.StringFlag{
				Name:    "api",
				Usage:   "specify the api backend, which will use /api as prefix",
				EnvVars: []string{"API"},
			},
			//
			&cli.StringSliceFlag{
				Name:    "header",
				Usage:   "specify the header, which will be added to the request",
				EnvVars: []string{"HEADERS"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		oauth2Provider := ctx.String("oauth2-provider")
		oauth2RedirectURI := ctx.String("oauth2-redirect-uri")
		if oauth2RedirectURI == "" {
			oauth2RedirectURI = fmt.Sprintf("http://127.0.0.1:%s/login/%s/callback", ctx.String("port"), oauth2Provider)
		}

		var headers http.Header
		if ctx.StringSlice("header") != nil {
			for _, header := range ctx.StringSlice("header") {
				headers = http.Header{}

				var kv []string
				if strings.Contains(header, "=") {
					kv = strings.SplitN(header, "=", 2)
				} else if strings.Contains(header, ":") {
					kv = strings.SplitN(header, ":", 2)
				} else {
					kv = []string{header}
				}

				if len(kv) == 2 {
					headers.Set(kv[0], kv[1])
				} else {
					headers.Set(kv[0], "")
				}
			}
		}

		return core.Serve(&core.Config{
			Port:                ctx.Int64("port"),
			Upstream:            ctx.String("upstream"),
			Prefix:              ctx.String("prefix"),
			DisableChangeOrigin: ctx.Bool("disable-change-origin"),
			//
			BasicUsername: ctx.String("basic-username"),
			BasicPassword: ctx.String("basic-password"),
			//
			BearerToken: ctx.String("bearer-token"),
			//
			AuthService: ctx.String("auth-service"),
			//
			Oauth2Provider:     oauth2Provider,
			Oauth2ClientID:     ctx.String("oauth2-client-id"),
			Oauth2ClientSecret: ctx.String("oauth2-client-secret"),
			Oauth2RedirectURI:  oauth2RedirectURI,
			//
			API: ctx.String("api"),
			//
			Headers: headers,
		})
	})

	app.Run()
}
