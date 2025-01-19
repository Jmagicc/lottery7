package router

import (
	"lottery7/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(lotteryHandler *handler.LotteryHandler) *gin.Engine {
	r := gin.Default()

	// API 路由组
	api := r.Group("/api")
	{
		api.GET("/lottery-results", lotteryHandler.GetLotteryResults)
		api.GET("/unused-numbers", lotteryHandler.GetUnusedNumbers)
		api.GET("/matrix", lotteryHandler.GetNumberMatrix)
		api.GET("/repeat-numbers", lotteryHandler.GetRepeatNumbers)
	}

	// 所有其他路由返回 index.html
	r.NoRoute(func(c *gin.Context) {
		c.File("index.html")
	})

	return r
}
