package auth

import "github.com/go-zoox/zoox"

func ApplyBearerToken(app *zoox.Application, bearerToken string) {
	if bearerToken == "" {
		return
	}

	app.Use(func(ctx *zoox.Context) {
		token, ok := ctx.BearerToken()
		if !ok {
			ctx.Status(401)
			return
		}

		if !(token == bearerToken) {
			ctx.Status(401)
			return
		}

		ctx.Next()
	})
}
