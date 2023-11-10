package core

import (
	"fmt"
	"net/http"

	"github.com/go-zoox/gzproxy/core/auth"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
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
	Oauth2Scope        string
	Oauth2BaseURL      string

	// api backend
	API string

	// proxy headers
	Headers http.Header
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	auth.ApplyBasicAuth(app, cfg.BasicUsername, cfg.BasicPassword)
	auth.ApplyBearerToken(app, cfg.BearerToken)
	auth.ApplyAuthService(app, cfg.AuthService)
	auth.ApplyOauth2(app, cfg.Oauth2Provider, cfg.Oauth2ClientID, cfg.Oauth2ClientSecret, cfg.Oauth2RedirectURI, cfg.Oauth2Scope, cfg.Oauth2BaseURL)

	if cfg.API != "" {
		app.Proxy("/api", cfg.API, func(sc *zoox.ProxyConfig) {
			if !cfg.DisableChangeOrigin {
				sc.ChangeOrigin = true
			}

			sc.Rewrites = rewriter.Rewriters{
				{
					From: "/api/(.*)",
					To:   "/$1",
				},
			}

			sc.RequestHeaders = cfg.Headers
		})
	}

	app.Proxy(".*", cfg.Upstream, func(sc *zoox.ProxyConfig) {
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

		sc.RequestHeaders = cfg.Headers
	})

	return app.Run(fmt.Sprintf(":%d", cfg.Port))
}
