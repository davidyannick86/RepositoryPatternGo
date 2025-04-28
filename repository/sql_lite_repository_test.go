package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSQLiteDatabase(t *testing.T) (*sql.DB, func()) {
	// Création d'une base de données SQLite en mémoire pour les tests
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Création du schéma de la base de données
	schema := `
		CREATE TABLE IF NOT EXISTS users (
		  id          TEXT PRIMARY KEY,
		  name        TEXT NOT NULL,
		  email       TEXT NOT NULL UNIQUE
		);`
	_, err = db.Exec(schema)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func setupSQLiteRepository(t *testing.T) (UserRepository, func()) {
	db, dbCleanup := setupSQLiteDatabase(t)

	// Création du repository
	repo := NewSqlLiteRepository(db)

	return repo, dbCleanup
}

func TestSqlLiteRepository_AddUser(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo, cleanup := setupSQLiteRepository(t)
	defer cleanup()

	// Données de test
	user := domain.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Exécution
	result, err := repo.AddUser(ctx, user)

	// Vérification
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)

	// Vérification que l'utilisateur a été ajouté à la base de données
	users, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)

	// Vérifier si notre nouvel utilisateur est dans la liste
	found := false
	for _, u := range users {
		if u.Email == user.Email {
			found = true
			break
		}
	}
	assert.True(t, found, "L'utilisateur devrait être trouvé dans la base de données")
}

func TestSqlLiteRepository_GetAllUsers(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo, cleanup := setupSQLiteRepository(t)
	defer cleanup()

	// Ajouter des utilisateurs de test
	testUsers := []domain.User{
		{Name: "User 1", Email: "user1@example.com"},
		{Name: "User 2", Email: "user2@example.com"},
		{Name: "User 3", Email: "user3@example.com"},
	}

	for _, user := range testUsers {
		_, err := repo.AddUser(ctx, user)
		require.NoError(t, err)
	}

	// Exécution
	users, err := repo.GetAllUsers(ctx)

	// Vérification
	require.NoError(t, err)
	assert.Equal(t, len(testUsers), len(users))

	// Vérifier que nos utilisateurs de test sont dans le résultat
	emails := make(map[string]bool)
	for _, user := range users {
		emails[user.Email] = true
	}

	for _, testUser := range testUsers {
		assert.True(t, emails[testUser.Email], "L'utilisateur avec l'email %s devrait être dans les résultats", testUser.Email)
	}
}

func TestSqlLiteRepository_GetAllUsers_Error(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo, cleanup := setupSQLiteRepository(t)

	// Force an error by closing the database connection before calling GetAllUsers
	cleanup() // This closes the database connection

	// Execute - this should fail because the connection is closed
	users, err := repo.GetAllUsers(ctx)

	// Verify we get an error
	assert.Error(t, err, "Should return an error when database connection is closed")
	assert.Nil(t, users, "Should not return users when there's an error")
}

func TestSqlLiteRepository_AddUser_Error(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo, cleanup := setupSQLiteRepository(t)

	// Test data
	user := domain.User{
		Name:  "Error Test User",
		Email: "error@example.com",
	}

	// Force an error by closing the database connection before calling AddUser
	cleanup() // This closes the database connection

	// Execute - this should fail because the connection is closed
	result, err := repo.AddUser(ctx, user)

	// Verify we get an error
	assert.Error(t, err, "Should return an error when database connection is closed")
	assert.Nil(t, result, "Should not return a user when there's an error")
}

func TestSqlLiteRepository_Integration(t *testing.T) {
	// Setup
	ctx := context.Background()
	repo, cleanup := setupSQLiteRepository(t)
	defer cleanup()

	// Test 1: Vérifier que la base de données est vide initialement
	initialUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	initialCount := len(initialUsers)
	assert.Equal(t, 0, initialCount, "La base de données devrait être vide initialement")

	// Test 2: Ajouter des utilisateurs initiaux
	initialTestUsers := []domain.User{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
	}

	for _, user := range initialTestUsers {
		_, err := repo.AddUser(ctx, user)
		require.NoError(t, err)
	}

	// Vérifier que les utilisateurs initiaux ont été ajoutés
	usersAfterInitial, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, len(initialTestUsers), len(usersAfterInitial))

	// Test 3: Ajouter un nouvel utilisateur
	newUser := domain.User{
		Name:  "Integration Test User",
		Email: "integration@example.com",
	}
	addedUser, err := repo.AddUser(ctx, newUser)
	require.NoError(t, err)
	assert.NotNil(t, addedUser)
	assert.NotEqual(t, uuid.Nil, addedUser.ID)

	// Test 4: Vérifier que l'utilisateur a été ajouté
	updatedUsers, err := repo.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, len(initialTestUsers)+1, len(updatedUsers))

	// Trouver notre nouvel utilisateur dans les résultats
	found := false
	for _, user := range updatedUsers {
		if user.Email == newUser.Email {
			found = true
			assert.Equal(t, newUser.Name, user.Name)
			break
		}
	}
	assert.True(t, found, "L'utilisateur ajouté devrait être trouvé dans la base de données")
}
