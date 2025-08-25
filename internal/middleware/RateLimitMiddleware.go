package middleware

import (
	"github.com/gin-gonic/gin"
	"go-chat/configs"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

var UserLimiter sync.Map
var InterfaceLimiter sync.Map

func UserRateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户的标识，有id则用id，无id则用ip
		id := c.GetString("id")
		if id == "" { // 没有id时使用ip
			id = c.ClientIP()
		}

		// 从 UserLimiter 中获取用户的限速器
		limiter, exists := UserLimiter.Load(id)
		if !exists {
			// 如果没有，初始化一个限速器（每秒最多10次请求，突发最大10次）
			limiter = rate.NewLimiter(rate.Every(1*time.Second), configs.AppConfig.Rate.UserLimit)
			UserLimiter.Store(id, limiter)
		}

		// 转换类型
		rateLimiter, ok := limiter.(*rate.Limiter)
		if !ok {
			// 错误处理，如果类型断言失败
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "内部错误，限速器转换失败。",
			})
			c.Abort() // 终止请求
			return
		}

		// 判断是否允许请求
		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试。",
			})
			c.Abort() // 终止请求
			return
		}

		// 如果没有超出限速，继续执行后续逻辑
		c.Next()
	}
}

func InterfaceRateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.Request.URL.Path
		method := c.Request.Method
		key := route + ":" + method
		limiter, exists := InterfaceLimiter.Load(key)
		if !exists {
			limiter = rate.NewLimiter(rate.Every(1*time.Second), configs.AppConfig.Rate.ApiLimit)
			InterfaceLimiter.Store(key, limiter)
		}
		rateLimiter, ok := limiter.(*rate.Limiter)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "内部错误，限速器转换失败。",
			})
			c.Abort() // 终止请求
			return
		}

		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试。",
			})
			c.Abort() // 终止请求
			return
		}
		c.Next()
	}
}
