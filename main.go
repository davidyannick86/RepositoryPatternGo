package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/davidyannick/repository-pattern/repository"
	service "github.com/davidyannick/repository-pattern/services"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	postgreSQL()
	sqlLite()
}

func postgreSQL() {
	dsn := "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Panicf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPsqlRepository(pool)
	ctx := context.Background()

	service := service.NewUserService(repo)

	user, err := service.AddUser(ctx, domain.User{
		Name:  generateRandomName(),  // #nosec G404
		Email: generateRandomEmail(), // #nosec G404
	})
	if err != nil {
		log.Panicf("Error adding user: %v", err)
	}

	log.Printf("Added user: %v", user)

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Panicf("Error getting users: %v", err)
	}
	log.Printf("All users: %v", users)
	log.Printf("Total users: %d", len(users))
}

func sqlLite() {
	db, err := sql.Open("sqlite3", "file:users.db?cache=shared&_fk=1")
	if err != nil {
		log.Panicf("ouvrir DB : %v", err)
	}
	defer db.Close()

	// On crée la table si nécessaire.
	schema := `
		CREATE TABLE IF NOT EXISTS users (
		  id          TEXT PRIMARY KEY,
		  name        TEXT NOT NULL,
		  email       TEXT NOT NULL UNIQUE
		);`
	if _, err2 := db.Exec(schema); err2 != nil {
		log.Panicf("création schema : %v", err2)
	}

	repo := repository.NewSQLLiteRepository(db)
	ctx := context.Background()

	service := service.NewUserService(repo)

	user, err := service.AddUser(ctx, domain.User{
		Name:  generateRandomName(),  // #nosec G404
		Email: generateRandomEmail(), // #nosec G404
	})
	if err != nil {
		log.Panicf("Error adding user: %v", err)
	}
	log.Printf("Added user: %v", user)

	users, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Panicf("Error getting users: %v", err)
	}
	log.Printf("All users: %v", users)
	log.Printf("Total users: %d", len(users))
}

func generateRandomEmail() string {
	// Use crypto/rand for cryptographically secure random numbers

	// Define possible characters for the username
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate random username of length 8-12
	n, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		log.Printf("Error generating random number: %v", err)
		n = big.NewInt(2) // Fallback
	}
	usernameLength := int(n.Int64()) + 8 // Random length between 8 and 12
	username := make([]byte, usernameLength)

	for i := range username {
		charIndex, randErr := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if randErr != nil {
			log.Printf("Error generating random index: %v", randErr)
			charIndex = big.NewInt(int64(i % len(charset))) // Fallback
		}
		username[i] = charset[charIndex.Int64()]
	}

	// Define list of possible domains
	domains := []string{
		"gmail.com",
		"yahoo.com",
		"hotmail.com",
		"outlook.com",
		"example.com",
	}

	// Choose a random domain
	domainIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(domains))))
	if err != nil {
		log.Printf("Error selecting random domain: %v", err)
		domainIndex = big.NewInt(0) // Fallback
	}
	domain := domains[domainIndex.Int64()]

	// Combine the parts to form an email
	return fmt.Sprintf("%s@%s", string(username), domain)
}

func generateRandomName() string {
	// Use crypto/rand for cryptographically secure random numbers

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
	characterIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(movieCharacters))))
	if err != nil {
		log.Printf("Error selecting random character: %v", err)
		characterIndex = big.NewInt(0) // Fallback
	}
	randomName := movieCharacters[characterIndex.Int64()]

	// Replace spaces with hyphens
	randomName = strings.ReplaceAll(randomName, " ", "-")

	// Generate random numbers for prefix and suffix
	num1, err := rand.Int(rand.Reader, big.NewInt(50))
	if err != nil {
		log.Printf("Error generating random number: %v", err)
		num1 = big.NewInt(1) // Fallback
	}

	num2, err := rand.Int(rand.Reader, big.NewInt(50))
	if err != nil {
		log.Printf("Error generating random number: %v", err)
		num2 = big.NewInt(2) // Fallback
	}

	randomName = fmt.Sprintf("%d-%s-%d", num1.Int64(), randomName, num2.Int64())
	return randomName
}
