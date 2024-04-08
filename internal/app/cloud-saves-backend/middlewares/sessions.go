package middlewares

import (
	sessions_store "cloud-saves-backend/internal/app/cloud-saves-backend/infra/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func Sessions(store sessions_store.Store) gin.HandlerFunc {
	store.SetOptions(&sessions.Options{
		HttpOnly: true,
		MaxAge:   86400,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/",
	})
	return sessions_store.Sessions("session", store)
}
