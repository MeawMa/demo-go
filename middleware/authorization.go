package middleware

import (
	"demo-go/models"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := ctx.Get("sub")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		e, _ := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		sub := user.(*models.User)  // the user that wants to access a resource.
		obj := ctx.Request.URL.Path // the resource that is going to be accessed.
		act := ctx.Request.Method   // the operation that the user performs on the resource.

		if res, _ := e.Enforce(sub, obj, act); res {
			// permit alice to read data1
			ctx.Next()
		} else {
			// deny the request, show an error
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}

	}
}
