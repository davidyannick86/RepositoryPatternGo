package repository

import (
	"context"
	"testing"
	"time"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongoContainer(t *testing.T) (testcontainers.Container, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(120 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return container, cleanup
}

func setupMongoRepository(t *testing.T, container testcontainers.Container) (UserRepository, func()) {
	ctx := context.Background()

	// Get host and port for connecting to the container
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "27017/tcp")
	require.NoError(t, err)

	// Create MongoDB client
	uri := "mongodb://" + host + ":" + port.Port()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	require.NoError(t, err)

	// Create repository
	repo := NewMongoDBRepository(client)

	// Initialize the collection by adding an index on email
	mongoClient := client.Database(databaseName).Collection(collectionName)
	_, err = mongoClient.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	require.NoError(t, err)

	cleanup := func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("failed to disconnect MongoDB client: %s", err)
		}
	}

	return repo, cleanup
}

func TestMongoDBRepository_AddUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupMongoContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupMongoRepository(t, container)
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

func TestMongoDBRepository_GetAllUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupMongoContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupMongoRepository(t, container)
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
	assert.Equal(t, len(testUsers), len(users))

	// Verify our test users are in the result
	emails := make(map[string]bool)
	for _, user := range users {
		emails[user.Email] = true
	}

	for _, testUser := range testUsers {
		assert.True(t, emails[testUser.Email], "User with email %s should be in the results", testUser.Email)
	}
}

func TestMongoDBRepository_Integration(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, containerCleanup := setupMongoContainer(t)
	defer containerCleanup()

	repo, repoCleanup := setupMongoRepository(t, container)
	defer repoCleanup()

	// Test 1: Get initial users (should be empty)
	initialUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	initialCount := len(initialUsers)
	assert.Equal(t, 0, initialCount, "Initial collection should be empty")

	// Test 2: Add initial test users
	initialTestUsers := []domain.User{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
	}

	for _, user := range initialTestUsers {
		_, err := repo.AddUser(ctx, user)
		require.NoError(t, err)
	}

	// Verify the initial users were added
	usersAfterInitial, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, len(initialTestUsers), len(usersAfterInitial))

	// Test 3: Add a new user
	newUser := domain.User{
		Name:  "Integration Test User",
		Email: "integration@example.com",
	}
	addedUser, err := repo.AddUser(ctx, newUser)
	require.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.NotEqual(t, uuid.Nil, addedUser.ID)

	// Test 4: Verify the user was added
	updatedUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, len(initialTestUsers)+1, len(updatedUsers))

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
