package app

import (
	"api/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weplanx/go/engine"
	"time"
)

func middleware(r *gin.Engine, values *common.Values) *gin.Engine {
	r.SetTrustedProxies(values.TrustedProxies)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     values.Cors.AllowOrigins,
		AllowMethods:     values.Cors.AllowMethods,
		AllowHeaders:     values.Cors.AllowHeaders,
		ExposeHeaders:    values.Cors.ExposeHeaders,
		AllowCredentials: values.Cors.AllowCredentials,
		MaxAge:           time.Duration(values.Cors.MaxAge) * time.Second,
	}))
	engine.RegisterValidation()
	return r
}

//func AuthGuard(passport *passport.Passport) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		tokenString := c.Cookies("access_token")
//		if tokenString == "" {
//			return c.JSON(fiber.Map{
//				"code":    401,
//				"message": common.LoginExpired.Error(),
//			})
//		}
//		claims, err := passport.Verify(tokenString)
//		if err != nil {
//			return err
//		}
//		c.Locals(common.TokenClaimsKey, claims)
//		return c.Next()
//	}
//}
