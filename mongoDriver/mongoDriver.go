package mongoDriver

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"quiz/model"
	"time"
)


func GetCollection(ctx context.Context,db , cl string) (*mongo.Collection,error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err!=nil {
		return nil,errors.Wrap(err,"Mongo connection error")
	}
	collection := client.Database(db).Collection(cl)
	return collection,nil
}

func GetHighScore(ctx context.Context,collection *mongo.Collection) []model.Score {

	findOptions := options.Find()
	findOptions.SetSort(bson.D{bson.E{Key: "total", Value: -1}, bson.E{Key: "level", Value: 1},bson.E{Key: "username", Value: 1}})
	findOptions.SetLimit(10)
	t := time.Now()
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())

	// apply findOptions
	cur, err := collection.Find(ctx, bson.M{"date":bson.M{"$gt":t1}}, findOptions)
	if err!=nil {
		log.Fatal(err)
	}
	var users []model.Score

	if err = cur.All(ctx,&users) ; err != nil {
		log.Fatal(err)
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cur, ctx)

	return users
}

func PutUser(ctx context.Context,collection *mongo.Collection,user model.User) (interface{},error) {

	res , er := collection.InsertOne(ctx,user)
	if er != nil {
		return nil,errors.Wrap(er,"Inserting error in PutUser function")
	}

	return res.InsertedID,nil
}


func GetUser(ctx context.Context,collection *mongo.Collection,username string) (model.User,error) {
	var user model.User
	if err := collection.FindOne(ctx, bson.M{"username":username}).Decode(&user); err != nil {
			return user,errors.Wrap(err,"User not found")
		}
	return user,nil
}


func GetUserData(ctx context.Context,collection *mongo.Collection,username string) (model.Score,error) {
	var user model.Score
	if err := collection.FindOne(ctx, bson.M{"username":username}).Decode(&user); err != nil {
		return user,errors.Wrap(err,"User not found")
	}
	return user,nil
}


func PutUserData(ctx context.Context,collection *mongo.Collection,userData model.Score)  (interface{},error){
	res , err := collection.InsertOne(ctx,userData)
	if err!=nil {
		return nil,errors.Wrap(err,"Error in inserting new user Data")
	}

	return res,nil
}

func Update(ctx context.Context,collection *mongo.Collection,user model.Score)  (interface{},error){
	update := bson.M{"$set": bson.M{"level":user.Level,"total":user.Total,"points":user.Points,"date":user.Date}}
	res , err := collection.UpdateOne(ctx,bson.M{"username":user.Username},update)
	if err != nil {
		return nil,errors.Wrap(err,"Storing Problem")
	}
	return res,nil
}


