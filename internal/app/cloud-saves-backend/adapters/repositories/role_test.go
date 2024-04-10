package repositories_test

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/repositories"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/tests"
	"context"
	"testing"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/stretchr/testify/require"
)

func TestRoleRepository (t *testing.T) {	
	db := tests.SetupSuite()

	roleRepo := repositories.NewRoleRepository(db, trmgorm.DefaultCtxGetter)

	ctx := context.Background()

	roles := []user.RoleName{user.RoleUser, user.RoleAdmin}

	t.Parallel()
	for _, roleName := range roles {
		roleName := roleName
		t.Run("GetRole USER", func(t *testing.T) {
			t.Parallel()

			role, err := roleRepo.GetByName(ctx, roleName)
			require.NoError(t, err)
			require.Equal(t, roleName, role.GetName())
		})
	}
}

