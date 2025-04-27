package repository

import (
	"context"

	"github.com/davidyannick/repository-pattern/domain"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoDBRepository struct {
	monClient *mongo.Client
}

var databaseName = "users"
var collectionName = "users"

func NewMongoDBRepository(client *mongo.Client) UserRepository {
	return &MongoDBRepository{monClient: client}
}

func (r *MongoDBRepository) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	user.ID = uuid.New()

	_, err := r.monClient.Database(databaseName).Collection(collectionName).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MongoDBRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.monClient.Database(databaseName).Collection(collectionName).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
