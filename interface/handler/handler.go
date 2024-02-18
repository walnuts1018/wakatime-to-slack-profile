package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/wakatime-to-slack-profile/config"
	"github.com/walnuts1018/wakatime-to-slack-profile/usecase"
)

var (
	uc *usecase.UpdateStatus
)

func NewHandler(cfg config.Config, updateStatus *usecase.UpdateStatus) (*gin.Engine, error) {
	uc = updateStatus

	r := gin.Default()
	store := cookie.NewStore([]byte(cfg.CookieSecret))
	r.Use(sessions.Sessions("WakatimeToSlack", store))
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/slack_signin")
	})
	r.GET("/signin", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/slack_signin")
	})

	r.GET("/slack_signin", slackSignIn)
	r.GET("/slack_callback", callback)

	return r, nil
}

func slackSignIn(ctx *gin.Context) {
	redirect := uc.SlackAuth()
	ctx.Redirect(http.StatusFound, redirect)
}

func wakatimeSignIn(ctx *gin.Context) {
	session := sessions.Default(ctx)
	state, redirect, err := uc.WakatimeAuth()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  fmt.Sprintf("failed to sign in: %v", err),
		})
		return
	}
	session.Set("state", state)
	session.Save()

	ctx.Redirect(http.StatusFound, redirect)
}

func callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	session := sessions.Default(ctx)
	if session.Get("state") != state {
		ctx.HTML(http.StatusBadRequest, "result.html", gin.H{
			"result": "error",
			"error":  "invalid state",
		})
		return
	}
	err := uc.Callback(ctx, code)
	if err != nil {
		slog.Error("failed to callback", "error", err)
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  "failed to callback",
		})
		return
	}

	err = uc.SetToken(ctx)
	if err != nil {
		slog.Error("failed to set token", "error", err)
		ctx.HTML(http.StatusInternalServerError, "result.html", gin.H{
			"result": "error",
			"error":  "failed to set token",
		})
		return
	}

	ctx.HTML(http.StatusOK, "result.html", gin.H{
		"result": "success",
	})
}
