package repositories_test

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/repositories"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/tests"
	"context"
	"fmt"
	"testing"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/stretchr/testify/require"
)

func TestTokenRepository(t *testing.T) {	
	db := tests.SetupSuite()

	userRepo := repositories.NewUserRepository(db, trmgorm.DefaultCtxGetter)
	tokenRepo := repositories.NewPasswordRecoveryTokenRepository(db, trmgorm.DefaultCtxGetter)

	ctx := context.Background()
	trManager := manager.Must(
		trmgorm.NewDefaultFactory(db),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested))),
		),
	)

	admin, err := userRepo.GetByUsername(ctx, "admin")  
	require.NoError(t, err)

	// trManager.Do(ctx, func(ctx context.Context) error {
	// 	token, err := tokenRepo.GetByUserId(ctx, admin.GetId())
	// 	require.NoError(t, err)
	// 	tokenRepo.Delete(ctx, token)
	// 	// require.Equal(t, admin.GetId(), token.GetUser().GetId())
	// 	// return fmt.Errorf("rollback")
	// 	return nil
	// })

	t.Run("No duplicates", func(t *testing.T) {
		tokens := []auth.PasswordRecoveryToken{
			*auth.NewPasswordRecoveryToken(admin),
			*auth.NewPasswordRecoveryToken(admin),
		}
	
		trManager.Do(ctx, func(ctx context.Context) error {
			token := tokens[0]
			err := tokenRepo.Create(ctx, &token)
			require.NoError(t, err)
	
			token = tokens[1]
			err = tokenRepo.Create(ctx, &token)
			require.Error(t, err)
			
			return fmt.Errorf("rollback")
		})
	})

	t.Run("GetTokenByUserId", func(t *testing.T) {
		tokens := []auth.PasswordRecoveryToken{
			*auth.NewPasswordRecoveryToken(admin),
		}
	
		trManager.Do(ctx, func(ctx context.Context) error {
			token := tokens[0]
			err := tokenRepo.Create(ctx, &token)
			require.NoError(t, err)
			
			tokenFromRepo, err := tokenRepo.GetByUserId(ctx, admin.GetId())
			require.NoError(t, err)
			require.Equal(t, admin.GetId(), tokenFromRepo.GetUser().GetId())
			return fmt.Errorf("rollback")
		})
	})

	t.Run("Delete", func(t *testing.T) {
		tokens := []auth.PasswordRecoveryToken{
			*auth.NewPasswordRecoveryToken(admin),
		}
	
		trManager.Do(ctx, func(ctx context.Context) error {
			token := tokens[0]
			err := tokenRepo.Create(ctx, &token)
			require.NoError(t, err)
	
			tokenFromRepo, err := tokenRepo.GetByUserId(ctx, admin.GetId())
			require.NoError(t, err)
			err = tokenRepo.Delete(ctx, tokenFromRepo)
			require.NoError(t, err)

			tokenFromRepo, err = tokenRepo.GetByUserId(ctx, admin.GetId())
			require.Error(t, err)
			require.Nil(t, tokenFromRepo)
			return fmt.Errorf("rollback")
		})
	})
}

