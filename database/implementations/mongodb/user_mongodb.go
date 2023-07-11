package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/model"
)

// UserMongoDB implements the UserDatabase interface using MongoDB
type UserMongoDB struct {
	db *mongo.Collection
}

// NewUserMongoDB creates a new UserMongoDB repository
func NewUserMongoDB(database *mongo.Database) *UserMongoDB {
	collection := database.Collection("users")
	return &UserMongoDB{
		db: collection,
	}
}

// GetUsers retrieves all users from MongoDB
func (m *UserMongoDB) GetUsers() ([]model.User, error) {
	var users []model.User

	cursor, err := m.db.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a user by ID from MongoDB
func (m *UserMongoDB) GetUserByID(id int) (model.User, error) {
	var user model.User

	filter := bson.M{"id": id}

	err := m.db.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// AddUser adds a new user to MongoDB
func (m *UserMongoDB) AddUser(user model.User) (model.User, error) {
	_, err := m.db.InsertOne(context.Background(), user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// UpdateUser updates a user in MongoDB
func (m *UserMongoDB) UpdateUser(user model.User) (model.User, error) {
	filter := bson.M{"id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"fullName": user.FullName,
			"userName": user.UserName,
			"email":    user.Email,
		},
	}

	_, err := m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// PatchUser partially updates a user in MongoDB
func (m *UserMongoDB) PatchUser(user model.User) (model.User, error) {
	filter := bson.M{"id": user.ID}

	update := bson.M{}

	if user.FullName != "" {
		update["$set"] = bson.M{"fullName": user.FullName}
	}
	if user.UserName != "" {
		update["$set"] = bson.M{"userName": user.UserName}
	}
	if user.Email != "" {
		update["$set"] = bson.M{"email": user.Email}
	}

	_, err := m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from MongoDB
func (m *UserMongoDB) DeleteUser(id int) error {
	filter := bson.M{"id": id}

	_, err := m.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
