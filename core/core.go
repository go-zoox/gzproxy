package core

import (
	"fmt"

	"github.com/go-zoox/gzproxy/core/auth"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox/defaults"
)

// Config is the basic config
type Config struct {
	Port int64

	// Upstream is the upstream service
	// Example: http://httpbin:8080
	Upstream string

	// Prefix is the prefix
	// Example: /v1
	Prefix string

	// BasicUsername is the basic username
	BasicUsername string
	// BasicPassword is the basic password
	BasicPassword string

	// mode: dynamic service with username and password

	// AuthService is auth service url
	// Example:
	//   POST https://example.com/api/login
	//	      Header => Content-Type: application/json
	//				Body => { "username": "username", "password": "password" }
	AuthService string
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	auth.ApplyServiceAuth(app, cfg.AuthService)
	auth.ApplyBasicAuth(app, cfg.BasicUsername, cfg.BasicPassword)

	app.Proxy(".*", cfg.Upstream, func(sc *proxy.SingleTargetConfig) {
		sc.ChangeOrigin = true

		if cfg.Prefix != "" {
			sc.Rewrites = rewriter.Rewriters{
				{
					From: fmt.Sprintf("%s/(.*)", cfg.Prefix),
					To:   "/$1",
				},
			}
		}
	})

	return app.Run()
}
