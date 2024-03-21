package auth

import (
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ExtractUser(ctx *gin.Context) (*userDTOs.UserResponseDTO, error) {
	session := sessions.Default(ctx)
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
