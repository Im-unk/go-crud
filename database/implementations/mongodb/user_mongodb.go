// package mongodb

// import (
// 	"context"
// 	"fmt"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"main.go/model"
// )

// // UserMongoDB implements the UserDatabase interface using MongoDB
// type UserMongoDB struct {
// 	db *mongo.Collection
// }

// // NewUserMongoDB creates a new UserMongoDB repository
// func NewUserMongoDB(database *mongo.Database) *UserMongoDB {
// 	collection := database.Collection("users")
// 	return &UserMongoDB{
// 		db: collection,
// 	}
// }

// // GetUsers retrieves all users from MongoDB
// func (m *UserMongoDB) GetUsers() ([]model.User, error) {
// 	var users []model.User

// 	cursor, err := m.db.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var user model.User
// 		err := cursor.Decode(&user)
// 		if err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// // GetUserByID retrieves a user by ID from MongoDB
// func (m *UserMongoDB) GetUserByID(id primitive.ObjectID) (model.User, error) {
// 	var user model.User

// 	// Print the ID before executing the query
// 	fmt.Println("Fetching user with ID:", id)

// 	filter := bson.M{"_id": id}

// 	err := m.db.FindOne(context.Background(), filter).Decode(&user)
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// // AddUser adds a new user to MongoDB
// func (m *UserMongoDB) AddUser(user model.User) (model.User, error) {
// 	_, err := m.db.InsertOne(context.Background(), user)
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// // UpdateUser updates a user in MongoDB
// func (m *UserMongoDB) UpdateUser(user model.User) (model.User, error) {
// 	filter := bson.M{"_id": user.ID}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"fullname": user.FullName,
// 			"username": user.UserName,
// 			"email":    user.Email,
// 		},
// 	}

// 	_, err := m.db.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// // PatchUser partially updates a user in MongoDB
// func (m *UserMongoDB) PatchUser(user model.User) (model.User, error) {
// 	filter := bson.M{"_id": user.ID}

// 	update := bson.M{}

// 	if user.FullName != "" {
// 		update["$set"] = bson.M{"fullname": user.FullName}
// 	}
// 	if user.UserName != "" {
// 		update["$set"] = bson.M{"username": user.UserName}
// 	}
// 	if user.Email != "" {
// 		update["$set"] = bson.M{"email": user.Email}
// 	}

// 	_, err := m.db.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		return model.User{}, err
// 	}

// 	return user, nil
// }

// func (m *UserMongoDB) DeleteUser(id primitive.ObjectID) error {
// 	filter := bson.M{"_id": id}

// 	_, err := m.db.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/cache"
	"main.go/model"
)

type UserMongoDB struct {
	db          *mongo.Collection
	cache       cache.Cacher
	cachePrefix string
}

func NewUserMongoDB(database *mongo.Database, cache cache.Cacher) *UserMongoDB {
	collection := database.Collection("users")
	return &UserMongoDB{
		db:          collection,
		cache:       cache,
		cachePrefix: "user:",
	}
}

func (m *UserMongoDB) GetUsers() ([]model.User, error) {
	var users []model.User
	cacheKey := "users"

	err := m.cache.Get(cacheKey, &users)
	if err != nil {
		// Cache miss, retrieve the users from the repository
		users, err = m.getUsersFromDB()
		if err != nil {
			return nil, err
		}

		// Store the users in the cache
		err = m.cache.Set(cacheKey, users, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set users in cache: %v\n", err)
		}
	}

	return users, nil
}

func (m *UserMongoDB) getUsersFromDB() ([]model.User, error) {
	var users []model.User

	cursor, err := m.db.Find(context.Background(), bson.M{})
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

func (m *UserMongoDB) GetUserByID(id primitive.ObjectID) (model.User, error) {
	var user model.User
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, id.Hex())

	err := m.cache.Get(cacheKey, &user)
	if err != nil {
		// Cache miss, retrieve the user from the repository
		user, err = m.getUserByIDFromDB(id)
		if err != nil {
			return model.User{}, err
		}

		// Store the user in the cache
		err = m.cache.Set(cacheKey, user, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set user in cache: %v\n", err)
		}
	}

	return user, nil
}

func (m *UserMongoDB) getUserByIDFromDB(id primitive.ObjectID) (model.User, error) {
	var user model.User
	filter := bson.M{"_id": id}

	err := m.db.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (m *UserMongoDB) GetLatestInsertedUser() (model.User, error) {
	// Sort the users by insertion time in descending order
	opts := options.FindOne().SetSort(bson.M{"_id": -1})

	var user model.User
	err := m.db.FindOne(context.Background(), bson.M{}, opts).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (m *UserMongoDB) AddUser(user model.User) (model.User, error) {
	addedUser, err := m.addUserToDB(user)
	if err != nil {
		return model.User{}, err
	}

	// Clear the users cache
	err = m.cache.Delete("users")
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete users cache: %v\n", err)
	}

	return addedUser, nil
}

func (m *UserMongoDB) addUserToDB(user model.User) (model.User, error) {
	_, err := m.db.InsertOne(context.Background(), user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (m *UserMongoDB) UpdateUser(user model.User) (model.User, error) {
	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"fullname": user.FullName,
			"username": user.UserName,
			"email":    user.Email,
		},
	}

	// Clear the user cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, user.ID.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	_, err = m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (m *UserMongoDB) PatchUser(user model.User) (model.User, error) {
	filter := bson.M{"_id": user.ID}

	update := bson.M{}

	if user.FullName != "" {
		update["$set"] = bson.M{"fullname": user.FullName}
	}
	if user.UserName != "" {
		update["$set"] = bson.M{"username": user.UserName}
	}
	if user.Email != "" {
		update["$set"] = bson.M{"email": user.Email}
	}

	// Clear the user cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, user.ID.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	_, err = m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (m *UserMongoDB) DeleteUser(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	// Clear the user cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, id.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	_, err = m.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
