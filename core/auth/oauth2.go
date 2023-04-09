package auth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/oauth2"
	oauth2Creator "github.com/go-zoox/oauth2/create"
	"github.com/go-zoox/random"
	"github.com/go-zoox/zoox"
)

func ApplyOauth2(app *zoox.Application, provider, clientID, clientSecret, redirectURI string) {
	if provider == "" || clientID == "" || clientSecret == "" || redirectURI == "" {
		return
	}

	client, err := oauth2Creator.Create(provider, &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create oauth2 client(provider: %s): %v", provider, err))
	}

	app.Use(func(ctx *zoox.Context) {
		if ctx.Method != "GET" {
			if ctx.Session().Get("oauth2.user") == "" {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401000,
						"message": "Unauthorized",
					})
					return
				}

				ctx.String(http.StatusUnauthorized, "Unauthorized")
				return
			}

			ctx.Next()
			return
		}

		// login
		if ctx.Path == "/login" || ctx.Path == fmt.Sprintf("/login/%s", provider) {
			originState := random.String(8)
			originFrom := ctx.Query().Get("from").String()
			client.Authorize(originState, func(loginUrl string) {
				if originFrom != "" {
					ctx.Session().Set("from", originFrom)
				}

				ctx.Session().Set("oauth2.state", originState)
				ctx.Redirect(loginUrl)
			})
			return
		}

		// callback
		if ctx.Path == "/login/callback" || ctx.Path == fmt.Sprintf("/login/%s/callback", provider) {
			code := ctx.Query().Get("code").String()
			state := ctx.Query().Get("state").String()

			originState := ctx.Session().Get("oauth2.state")
			if state != originState {
				logger.Errorf("invalid oauth2 state, expect %s, but got %s", originState, state)
				time.Sleep(1 * time.Second)
				ctx.Redirect(fmt.Sprintf("/login?reason=%s", "invalid_oauth2_state"))
				return
			}
			originFrom := ctx.Session().Get("from")
			if originFrom == "" {
				originFrom = "/"
			}

			client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
				userSessionKey := fmt.Sprintf("user:%s", user.ID)

				ctx.Cache().Set(userSessionKey, user, ctx.App.SessionMaxAge)

				ctx.Session().Set("oauth2.user", userSessionKey)
				// ctx.Session().Set("oauth2.token", token.AccessToken)

				ctx.Redirect(originFrom)
			})
			return
		}

		// logout
		if ctx.Path == "/logout" {
			client.Logout(func(logoutUrl string) {
				ctx.Session().Del("oauth2.user")
				ctx.Redirect(logoutUrl)
			})
			return
		}

		if ctx.Path == "/api/user" {
			userSessionKey := ctx.Session().Get("oauth2.user")
			if userSessionKey == "" {
				if ctx.AcceptJSON() {
					ctx.JSON(http.StatusUnauthorized, zoox.H{
						"code":    401001,
						"message": "unauthorized",
					})
					return
				}

				ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Path), "user not login or token expired"))
				return
			}
			user := oauth2.User{}
			if err := ctx.Cache().Get(userSessionKey, &user); err != nil {
				time.Sleep(1 * time.Second)
				ctx.Redirect(fmt.Sprintf("/login?from=%s&reason=%s", url.QueryEscape(ctx.Path), "user cache not found"))
				return
			}

			ctx.JSON(http.StatusOK, user)
			return
		}

		if ctx.Session().Get("oauth2.user") == "" {
			originFrom := ctx.Request.RequestURI
			time.Sleep(1 * time.Second)
			ctx.Redirect(fmt.Sprintf("/login?from=%s", url.QueryEscape(originFrom)))
			return
		}

		userSessionKey := ctx.Session().Get("oauth2.user")
		user := oauth2.User{}
		if err := ctx.Cache().Get(userSessionKey, &user); err != nil {
			time.Sleep(1 * time.Second)
			ctx.Redirect(fmt.Sprintf("/login?reason=%s", "user cache not found"))
			return
		}

		token, err := ctx.Jwt().Sign(map[string]interface{}{
			"user_id":       user.ID,
			"user_nickname": user.Nickname,
			"user_avatar":   user.Avatar,
			"user_email":    user.Email,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Request.Header.Set("X-GZAuth-Token", token)

		ctx.Next()
	})
}
