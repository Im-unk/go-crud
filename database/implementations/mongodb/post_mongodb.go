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

type PostMongoDB struct {
	db          *mongo.Collection
	cache       cache.Cacher
	cachePrefix string
}

func NewPostMongoDB(database *mongo.Database, cache cache.Cacher) *PostMongoDB {
	collection := database.Collection("posts")
	return &PostMongoDB{
		db:          collection,
		cache:       cache,
		cachePrefix: "post:",
	}
}

func (m *PostMongoDB) GetPosts() ([]model.Post, error) {
	var posts []model.Post
	cacheKey := "posts"

	err := m.cache.Get(cacheKey, &posts)
	if err != nil {
		// Cache miss, retrieve the posts from the repository
		posts, err = m.getPostsFromDB()
		if err != nil {
			return nil, err
		}

		// Store the posts in the cache
		err = m.cache.Set(cacheKey, posts, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set posts in cache: %v\n", err)
		}
	}

	return posts, nil
}

func (m *PostMongoDB) getPostsFromDB() ([]model.Post, error) {
	var posts []model.Post

	cursor, err := m.db.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var post model.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *PostMongoDB) GetPostByID(id primitive.ObjectID) (model.Post, error) {
	var post model.Post
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, id.Hex())

	err := m.cache.Get(cacheKey, &post)
	if err != nil {
		// Cache miss, retrieve the post from the repository
		post, err = m.getPostByIDFromDB(id)
		if err != nil {
			return model.Post{}, err
		}

		// Store the post in the cache
		err = m.cache.Set(cacheKey, post, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set post in cache: %v\n", err)
		}
	}

	return post, nil
}

func (m *PostMongoDB) getPostByIDFromDB(id primitive.ObjectID) (model.Post, error) {
	var post model.Post
	filter := bson.M{"_id": id}

	err := m.db.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) GetLatestInsertedPost() (model.Post, error) {
	// Sort the posts by insertion time in descending order
	opts := options.FindOne().SetSort(bson.M{"_id": -1})

	var post model.Post
	err := m.db.FindOne(context.Background(), bson.M{}, opts).Decode(&post)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) AddPost(post model.Post) (model.Post, error) {
	addedPost, err := m.addPostToDB(post)
	if err != nil {
		return model.Post{}, err
	}

	// Clear the posts cache
	err = m.cache.Delete("posts")
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete posts cache: %v\n", err)
	}

	return addedPost, nil
}

func (m *PostMongoDB) addPostToDB(post model.Post) (model.Post, error) {
	_, err := m.db.InsertOne(context.Background(), post)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) UpdatePost(post model.Post) (model.Post, error) {
	filter := bson.M{"_id": post.ID}

	update := bson.M{
		"$set": bson.M{
			"title": post.Title,
			"body":  post.Body,
		},
	}

	// Clear the post cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, post.ID.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	_, err = m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) PatchPost(post model.Post) (model.Post, error) {
	filter := bson.M{"_id": post.ID}

	update := bson.M{}

	if post.Title != "" {
		update["$set"] = bson.M{"title": post.Title}
	}
	if post.Body != "" {
		update["$set"] = bson.M{"body": post.Body}
	}

	// Clear the post cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, post.ID.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	_, err = m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) DeletePost(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	// Clear the post cache
	cacheKey := fmt.Sprintf("%s%s", m.cachePrefix, id.Hex())
	err := m.cache.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	_, err = m.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
