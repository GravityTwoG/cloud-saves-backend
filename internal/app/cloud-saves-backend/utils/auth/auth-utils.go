package auth

import (
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	sessions_store "cloud-saves-backend/internal/app/cloud-saves-backend/sessions"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ExtractUser(ctx *gin.Context) (*userDTOs.UserResponseDTO, error) {
	session := sessions_store.Default(ctx)
	user := session.Get("user")
	if user == nil {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	userResponseDTO, ok := user.(*userDTOs.UserResponseDTO)
	if !ok {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	return userResponseDTO, nil
}
