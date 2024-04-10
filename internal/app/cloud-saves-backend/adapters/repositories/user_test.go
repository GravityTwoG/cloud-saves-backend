package repositories_test

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/repositories"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
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

func TestUserRepository(t *testing.T) {	
	db := tests.SetupSuite()

	userRepo := repositories.NewUserRepository(db, trmgorm.DefaultCtxGetter)
	roleRepo := repositories.NewRoleRepository(db, trmgorm.DefaultCtxGetter)

	ctx := context.Background()
	trManager := manager.Must(
		trmgorm.NewDefaultFactory(db),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested))),
		),
	)

	roleUser, err := roleRepo.GetByName(ctx, user.RoleUser)
	require.NoError(t, err)
	roleAdmin, err := roleRepo.GetByName(ctx, user.RoleAdmin)
	require.NoError(t, err)

	t.Parallel()
	t.Run("No duplicate usernames", func(t *testing.T) {
		t.Parallel()
		user1, err := user.NewUser(
			"username_duplicate", "username_duplicate1@mail.ru", 
			"password", roleUser,
		)
		require.NoError(t, err)
		user2, err := user.NewUser(
			"username_duplicate", "username_duplicate2@mail.ru", 
			"password", roleUser,
		) 
		require.NoError(t, err)
	
		trManager.Do(ctx, func(ctx context.Context) error {
			err := userRepo.Create(ctx, user1)
			require.NoError(t, err)

			err = userRepo.Create(ctx, user2)
			require.Error(t, err)
			
			return fmt.Errorf("rollback")
		})
	})

	t.Run("No duplicate emails", func(t *testing.T) {
		t.Parallel()
		user1, err := user.NewUser(
			"email_duplicate1", "email_duplicate@mail.ru", 
			"password", roleUser,
		)
		require.NoError(t, err)
		user2, err := user.NewUser(
			"email_duplicate2", "email_duplicate@mail.ru", 
			"password", roleUser,
		) 
		require.NoError(t, err)
	
		trManager.Do(ctx, func(ctx context.Context) error {
			err := userRepo.Create(ctx, user1)
			require.NoError(t, err)

			err = userRepo.Create(ctx, user2)
			require.Error(t, err)
			
			return fmt.Errorf("rollback")
		})
	})

	t.Run("GetUserById", func(t *testing.T) {
		t.Parallel()
		user, err := user.NewUser(
			"username", "username1@mail.ru", 
			"password", roleAdmin,
		)
		require.NoError(t, err) 
	
		trManager.Do(ctx, func(ctx context.Context) error {
			err := userRepo.Create(ctx, user)
			require.NoError(t, err)
			
			userFromRepo, err := userRepo.GetById(ctx, user.GetId())
			require.NoError(t, err)
			require.Equal(t, user.GetId(), userFromRepo.GetId())
			return fmt.Errorf("rollback")
		})
	})

	t.Run("GetByUsername", func(t *testing.T) {
		t.Parallel()
		user, err := user.NewUser(
			"username", "username1@mail.ru", 
			"password", roleAdmin,
		)
		require.NoError(t, err) 
	
		trManager.Do(ctx, func(ctx context.Context) error {
			err := userRepo.Create(ctx, user)
			require.NoError(t, err)
			
			userFromRepo, err := userRepo.GetByUsername(ctx, user.GetUsername())
			require.NoError(t, err)
			require.Equal(t, user.GetId(), userFromRepo.GetId())
			require.Equal(t, user.GetUsername(), userFromRepo.GetUsername())
			return fmt.Errorf("rollback")
		})
	})

	t.Run("GetByEmail", func(t *testing.T) {
		t.Parallel()
		user, err := user.NewUser(
			"username_get_email", "username_email_get@mail.ru", 
			"password", roleAdmin,
		)
		require.NoError(t, err) 
	
		trManager.Do(ctx, func(ctx context.Context) error {
			err := userRepo.Create(ctx, user)
			require.NoError(t, err)
			
			userFromRepo, err := userRepo.GetByEmail(ctx, user.GetEmail())
			require.NoError(t, err)
			require.Equal(t, user.GetId(), userFromRepo.GetId())
			require.Equal(t, user.GetEmail(), userFromRepo.GetEmail())
			return fmt.Errorf("rollback")
		})
	})
}

