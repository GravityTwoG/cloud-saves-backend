package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func Sessions(redisHost string, sessionSecret []byte) gin.HandlerFunc {
	store, err := redis.NewStore(10, "tcp", redisHost, "", sessionSecret)
	if err != nil {
		panic(err)
	}

	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})
	return sessions.Sessions("session", store)
}
