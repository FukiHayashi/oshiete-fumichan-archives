package router

import (
	"net/http"
	"takanome/controllers"
	"takanome/render"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// HTMLテンプレートパス設定
	router.HTMLRender = render.CreateRenderTemplates()

	// 静的パス設定
	router.Static("./assets", "./assets")

	root := router.Group("/")
	{
		root.GET("/", func(ctx *gin.Context) { ctx.HTML(http.StatusOK, "index.html", nil) })
		root.GET("/tweets", controllers.TweetsHandler)
		root.GET("/tweets/:tag", controllers.TweetsTagHandler)
		root.GET("/jobrunner/status", func(c *gin.Context) { c.JSON(http.StatusOK, jobrunner.StatusJson()) })
	}
	return router
}
