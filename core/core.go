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

	// DisableChangeOrigin is disable change origin
	DisableChangeOrigin bool

	// BasicUsername is the basic username
	BasicUsername string
	// BasicPassword is the basic password
	BasicPassword string

	// BearerToken is the bearer token
	BearerToken string

	// mode: dynamic service with username and password

	// AuthService is auth service url
	// Example:
	//   POST https://example.com/api/login
	//	      Header => Content-Type: application/json
	//				Body => { "username": "username", "password": "password" }
	AuthService string

	// oauth2
	Oauth2Provider     string
	Oauth2ClientID     string
	Oauth2ClientSecret string
	Oauth2RedirectURI  string
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	auth.ApplyBasicAuth(app, cfg.BasicUsername, cfg.BasicPassword)
	auth.ApplyBearerToken(app, cfg.BearerToken)
	auth.ApplyAuthService(app, cfg.AuthService)
	auth.ApplyOauth2(app, cfg.Oauth2Provider, cfg.Oauth2ClientID, cfg.Oauth2ClientSecret, cfg.Oauth2RedirectURI)

	app.Proxy(".*", cfg.Upstream, func(sc *proxy.SingleTargetConfig) {
		if !cfg.DisableChangeOrigin {
			sc.ChangeOrigin = true
		}

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
