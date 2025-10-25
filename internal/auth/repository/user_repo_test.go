// internal/auth/repository/user_repository_integration_test.go
package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"swasthAI/internal/auth/models"
	"swasthAI/internal/auth/repository"
	"swasthAI/pkg/domain_errors"
	"swasthAI/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	testDB      *bun.DB
	pgContainer *postgres.PostgresContainer
	cleanupOnce func()
	testLogger  logger.Logger
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// === Start PostgreSQL container (NEW API) ===
	pg, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	pgContainer = pg

	// === Get connection string ===
	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	// === Connect via pgdriver (recommended for Bun) ===
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqlDB := sql.OpenDB(connector)
	testDB = bun.NewDB(sqlDB, pgdialect.New())

	// === Enable uuid-ossp ===
	_, err = testDB.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		log.Fatalf("failed to create extension: %v", err)
	}

	// === Run migrations ===
	_, err = testDB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	// === Run tests ===
	code := m.Run()

	// === Cleanup ===
	testDB.Close()
	_ = pgContainer.Terminate(ctx)

	os.Exit(code)
}
func TestUserRepository_Integration(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewUserRepository(testDB, testLogger)

	// === Helper: Reset DB before each subtest ===
	resetDB := func(t *testing.T) {
		_, err := testDB.ExecContext(ctx, `TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
		require.NoError(t, err)
	}

	t.Run("Create and FindByPhone", func(t *testing.T) {
		resetDB(t)

		user := &models.User{
			Phone:     "+919876543210",
			FirstName: "रमेश",
			LastName:  "कुमार",
			Language:  "hi",
		}
		user.PrepareCreate()

		created, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "+919876543210", created.Phone)

		found, err := repo.FindByPhone(ctx, "+919876543210")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "रमेश", found.FirstName)
	})

	t.Run("FindByPhone - Not Found", func(t *testing.T) {
		resetDB(t)

		found, err := repo.FindByPhone(ctx, "+919999999999")
		require.Error(t, err)                                         // Expect an error
		assert.Nil(t, found)                                          // User should be nil
		assert.True(t, errors.Is(err, domain_errors.ErrUserNotFound)) // Specific error check
	})
	// t.Run("FindByPhone - Not Found", func(t *testing.T) {
	// 	resetDB(t)

	// 	found, err := repo.FindByPhone(ctx, "+919999999999")
	// 	require.NoError(t, err)
	// 	assert.Nil(t, found) // Should be nil, not error
	// })

	t.Run("FindByID", func(t *testing.T) {
		resetDB(t)

		user := &models.User{
			Phone:     "+918888888888",
			FirstName: "सूरज",
			LastName:  "शर्मा",
		}
		user.PrepareCreate()
		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
	})

	t.Run("Update", func(t *testing.T) {
		resetDB(t)

		user := &models.User{
			Phone:     "+917777777777",
			FirstName: "अमित",
			LastName:  "वर्मा",
			Language:  "en",
		}
		user.PrepareCreate()
		created, err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Modify fields
		created.FirstName = "अमिताभ"
		created.Language = "ta"
		created.FullName = "अमिताभ वर्मा" // Must set manually for Bun

		updated, err := repo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, "अमिताभ", updated.FirstName)
		assert.Equal(t, "ta", updated.Language)
		assert.Equal(t, "अमिताभ वर्मा", updated.FullName)
	})

	t.Run("Create - Duplicate Phone", func(t *testing.T) {
		resetDB(t)

		user1 := &models.User{
			Phone:     "+915555555555",
			FirstName: "Test",
			LastName:  "User",
		}
		user1.PrepareCreate()
		_, err := repo.Create(ctx, user1)
		require.NoError(t, err)

		user2 := &models.User{
			Phone:     "+915555555555",
			FirstName: "Duplicate",
			LastName:  "User",
		}
		user2.PrepareCreate()
		_, err = repo.Create(ctx, user2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
	})
}
