package auth

import (
	sessions_store "cloud-saves-backend/internal/app/cloud-saves-backend/adapters/sessions"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/ports/dto/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ExtractUser(ctx *gin.Context) (*userDTOs.UserDTO, error) {
	session := sessions_store.Default(ctx)
	user := session.Get("user")
	if user == nil {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	userResponseDTO, ok := user.(*userDTOs.UserDTO)
	if !ok {
		return nil, fmt.Errorf("UNAUTHORIZED")
	}

	return userResponseDTO, nil
}
