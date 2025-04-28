package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgresContainer(t *testing.T) (*postgres.PostgresContainer, func()) {
	ctx := context.Background()

	// Get the absolute path to init.sql
	workDir, err := os.Getwd()
	require.NoError(t, err)

	rootDir := filepath.Dir(workDir)
	initScriptPath := filepath.Join(rootDir, "init.sql")

	// Create and start the PostgreSQL container with increased timeout
	container, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second)),
	)

	require.NoError(t, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return container, cleanup
}

func setupRepository(t *testing.T, container *postgres.PostgresContainer) (UserRepository, func()) {
	ctx := context.Background()

	// Get connection details
	connString, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	// Create connection pool
	pool, err := pgxpool.New(ctx, connString)
	require.NoError(t, err)

	// Create repository
	repo := NewPsqlRepository(pool)

	cleanup := func() {
		pool.Close()
	}

	return repo, cleanup
}

func TestPsqlRepository_AddUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupPostgresContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupRepository(t, container)
	defer repoCleanup()

	// Test data
	user := domain.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Execute
	result, err := repo.AddUser(ctx, user)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)

	// Verify it was added to the database
	users, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)

	// Check if our new user is in the list
	found := false
	for _, u := range users {
		if u.Email == user.Email {
			found = true
			break
		}
	}
	assert.True(t, found, "User should be found in the database")
}

func TestPsqlRepository_GetAllUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupPostgresContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupRepository(t, container)
	defer repoCleanup()

	// Add some test users
	testUsers := []domain.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range testUsers {
		_, err := repo.AddUser(ctx, user)
		require.NoError(t, err)
	}

	// Execute
	users, err := repo.GetAllUsers(ctx)

	// Verify
	require.NoError(t, err)

	// We should have at least the users we added plus any from init.sql
	assert.GreaterOrEqual(t, len(users), len(testUsers))

	// Verify our test users are in the result
	emails := make(map[string]bool)
	for _, user := range users {
		emails[user.Email] = true
	}

	for _, testUser := range testUsers {
		assert.True(t, emails[testUser.Email], fmt.Sprintf("User with email %s should be in the results", testUser.Email))
	}
}

func TestPsqlRepository_GetAllUsers_Error(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupPostgresContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupRepository(t, container)

	// Force an error by closing the connection pool before calling GetAllUsers
	repoCleanup() // This will close the pool

	// Execute - this should fail because the connection is closed
	users, err := repo.GetAllUsers(ctx)

	// Verify we get an error
	assert.Error(t, err, "Should return an error when database connection is closed")
	assert.Nil(t, users, "Should not return users when there's an error")
}

func TestPsqlRepository_AddUser_Error(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupPostgresContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupRepository(t, container)

	// Test data
	user := domain.User{
		Name:  "Error Test User",
		Email: "error@example.com",
	}

	// Force an error by closing the connection pool before calling AddUser
	repoCleanup() // This will close the pool

	// Execute - this should fail because the connection is closed
	result, err := repo.AddUser(ctx, user)

	// Verify we get an error
	assert.Error(t, err, "Should return an error when database connection is closed")
	assert.Nil(t, result, "Should not return a user when there's an error")
}

func TestPsqlRepository_Integration(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupPostgresContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupRepository(t, container)
	defer repoCleanup()

	// Test 1: Get initial users from init.sql
	initialUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	initialCount := len(initialUsers)

	// The init.sql script should have added 2 users
	assert.GreaterOrEqual(t, initialCount, 2, "Should have at least 2 users from init.sql")

	// Test 2: Add a new user
	newUser := domain.User{
		Name:  "Integration Test User",
		Email: "integration@example.com",
	}
	addedUser, err := repo.AddUser(ctx, newUser)
	require.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.NotEqual(t, uuid.Nil, addedUser.ID)

	// Test 3: Verify the user was added
	updatedUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, initialCount+1, len(updatedUsers))

	// Find our new user in the results
	found := false
	for _, user := range updatedUsers {
		if user.Email == newUser.Email {
			found = true
			assert.Equal(t, newUser.Name, user.Name)
			break
		}
	}
	assert.True(t, found, "Added user should be found in the database")
}
