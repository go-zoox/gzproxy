package auth

import (
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
)

func ApplyAuthService(app *zoox.Application, authService string) {
	if authService == "" {
		return
	}

	app.Use(func(ctx *zoox.Context) {
		user, pass, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.Set("WWW-Authenticate", `Basic realm="go-zoox"`)
			ctx.Status(401)
			return
		}

		response, err := fetch.Post(authService, &fetch.Config{
			Headers: fetch.Headers{
				"Content-Type": "application/json",
			},
			Body: map[string]string{
				"from":     "go-zoox/gzproxy.basic",
				"username": user,
				"password": pass,
			},
		})
		if err != nil {
			logger.Errorf("basic auth with auth-service error: %s", err)
			fmt.PrintJSON(map[string]any{
				"request":  response.Request,
				"response": response.String(),
			})

			ctx.String(500, "internal server error")
			return
		}

		if response.Status != 200 {
			ctx.String(400, "invalid username and password: %s", response.String())
			return
		}

		ctx.Next()
	})
}
