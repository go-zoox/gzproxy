package auth

import "github.com/go-zoox/zoox"

func ApplyBasicAuth(app *zoox.Application, username, password string) {
	if username == "" || password == "" {
		return
	}

	app.Use(func(ctx *zoox.Context) {
		user, pass, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.Set("WWW-Authenticate", `Basic realm="go-zoox"`)
			ctx.Status(401)
			return
		}

		if !(user == username && pass == password) {
			ctx.Status(401)
			return
		}

		ctx.Next()
	})
}
