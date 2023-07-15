package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main.go/model"
)

// PostMongoDB implements the PostDatabase interface using MongoDB
type PostMongoDB struct {
	db *mongo.Collection
}

// NewPostMongoDB creates a new PostMongoDB repository
func NewPostMongoDB(database *mongo.Database) *PostMongoDB {
	collection := database.Collection("posts")
	return &PostMongoDB{
		db: collection,
	}
}

// GetPosts retrieves all posts from MongoDB
func (m *PostMongoDB) GetPosts() ([]model.Post, error) {
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

// GetPostByID retrieves a post by ID from MongoDB
func (m *PostMongoDB) GetPostByID(id int) (model.Post, error) {
	var post model.Post

	filter := bson.M{"_id": id}

	err := m.db.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

// AddPost adds a new post to MongoDB
func (m *PostMongoDB) AddPost(post model.Post) (model.Post, error) {
	_, err := m.db.InsertOne(context.Background(), post)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

func (m *PostMongoDB) UpdatePost(post model.Post) (model.Post, error) {
	filter := bson.M{"_id": post.ID} // Use "_id" instead of "id"

	update := bson.M{
		"$set": bson.M{
			"title": post.Title,
			"body":  post.Body,
		},
	}

	_, err := m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

// PatchPost partially updates a post in MongoDB
func (m *PostMongoDB) PatchPost(post model.Post) (model.Post, error) {
	filter := bson.M{"_id": post.ID}

	update := bson.M{}

	if post.Title != "" {
		update["$set"] = bson.M{"title": post.Title}
	}
	if post.Body != "" {
		update["$set"] = bson.M{"body": post.Body}
	}

	_, err := m.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model.Post{}, err
	}

	return post, nil
}

// DeletePost deletes a post by ID from MongoDB
func (m *PostMongoDB) DeletePost(id int) error {
	filter := bson.M{"_id": id}

	_, err := m.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
