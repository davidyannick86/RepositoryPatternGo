package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/davidyannick/repository-pattern/repository"
	service "github.com/davidyannick/repository-pattern/services"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	postGreSql()
	sqlLite()
}

func postGreSql() {
	dsn := "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPsqlRepository(pool)
	ctx := context.Background()

	service := service.NewUserService(repo)

	user, err := service.AddUser(ctx, domain.User{
		Name:  generateRandomName(),
		Email: generateRandomEmail(),
	})
	if err != nil {
		log.Fatalf("Error adding user: %v", err)
	}
	log.Printf("Added user: %v", user)

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}
	log.Printf("All users: %v", users)
	log.Printf("Total users: %d", len(users))
}

func sqlLite() {
	db, err := sql.Open("sqlite3", "file:users.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("ouvrir DB : %v", err)
	}
	defer db.Close()

	// On crée la table si nécessaire.
	schema := `
		CREATE TABLE IF NOT EXISTS users (
		  id          TEXT PRIMARY KEY,
		  name        TEXT NOT NULL,
		  email       TEXT NOT NULL UNIQUE
		);`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("création schema : %v", err)
	}

	repo := repository.NewSQLLiteRepository(db)
	ctx := context.Background()

	service := service.NewUserService(repo)

	user, err := service.AddUser(ctx, domain.User{
		Name:  generateRandomName(),
		Email: generateRandomEmail(),
	})
	if err != nil {
		log.Fatalf("Error adding user: %v", err)
	}
	log.Printf("Added user: %v", user)

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}
	log.Printf("All users: %v", users)
	log.Printf("Total users: %d", len(users))
}

func generateRandomEmail() string {
	// Initialize random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define possible characters for the username
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate random username of length 8-12
	usernameLength := r.Intn(5) + 8 // Random length between 8 and 12
	username := make([]byte, usernameLength)
	for i := range username {
		username[i] = charset[r.Intn(len(charset))]
	}

	// Define list of possible domains
	domains := []string{
		"gmail.com",
		"yahoo.com",
		"hotmail.com",
		"outlook.com",
		"example.com",
		"mail.com",
	}

	// Select a random domain
	domain := domains[r.Intn(len(domains))]

	// Combine username and domain to form email
	return generateRandomName() + "@" + domain
}

func generateRandomName() string {
	// Initialize random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define a list of movie character names
	movieCharacters := []string{
		"James Bond",
		"Indiana Jones",
		"Ellen Ripley",
		"Luke Skywalker",
		"Darth Vader",
		"Tony Stark",
		"Bruce Wayne",
		"Hermione Granger",
		"Harry Potter",
		"Forrest Gump",
		"Michael Corleone",
		"Rick Blaine",
		"Jack Sparrow",
		"Hannibal Lecter",
		"Dorothy Gale",
		"Rocky Balboa",
		"Sarah Connor",
		"Marty McFly",
		"Han Solo",
		"Katniss Everdeen",
		"John McClane",
		"Vito Corleone",
		"Norman Bates",
	}

	// Choose a random character
	randomName := movieCharacters[r.Intn(len(movieCharacters))]
	// Replace spaces with hyphens
	randomName = strings.ReplaceAll(randomName, " ", "-")
	randomName = fmt.Sprintf("%d-%s-%d", rand.Intn(50), randomName, rand.Intn(50))
	return randomName
}
