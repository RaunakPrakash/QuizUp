package userController

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"quiz/model"
	"quiz/mongoDriver"
)


type Mongo struct {
	collection *mongo.Collection
	user model.User
	userData model.Score
}


func (m *Mongo) SetCollection(ctx context.Context,db,cl string) error {
	col,err := mongoDriver.GetCollection(ctx,db,cl)
	if err!=nil {
		return errors.Wrap(err,"Collection issue in SetCollection method in UserController package")
	}
	m.collection = col
	return nil
}

func (m *Mongo) SetUser(user model.User)  {
	m.user = user
}

func (m *Mongo) SetUserData(user model.Score)  {
	m.userData = user
}


func (m *Mongo) Get(ctx context.Context,db,cl,username string)  (model.User,bool){
	err:= m.SetCollection(ctx,db,cl)
	if err != nil {
		return model.User{},false
	}
	user,er:=mongoDriver.GetUser(ctx, m.collection,username)
	if er != nil{
		return user,false
	}
	return user,true
}

func (m *Mongo) GetUserData(ctx context.Context,username string)  (model.Score,error){
	user,er:=mongoDriver.GetUserData(ctx, m.collection,username)
	if er!=nil{
		return model.Score{},errors.Wrap(er,"Issue in Fetching User Data in GetUserData method in UserController")
	}
	return user,nil
}

func (m *Mongo) GetHighScore(ctx context.Context)  []model.Score{
	user  := mongoDriver.GetHighScore(ctx, m.collection)
	return user
}



func (m *Mongo) Put(ctx context.Context)  (interface{},error){
	password, err := m.HashPassword(m.user.Password)
	if err != nil {
		return nil, err
	}
	m.user.Password = password
	putUser, err := mongoDriver.PutUser(ctx, m.collection, m.user)
	if err != nil {
		return nil,errors.Wrap(err,"putting user")
	}
	return putUser,nil
}

func (m *Mongo) PutUserData(ctx context.Context)  (interface{},error){
	putUser, err := mongoDriver.PutUserData(ctx, m.collection, m.userData)
	if err != nil {
		return nil,errors.Wrap(err,"putting user")
	}
	return putUser,nil
}

func (m *Mongo) Reset(ctx context.Context)  (interface{},error){
	update, err := mongoDriver.Update(ctx, m.collection, m.userData)
	if err != nil {
		return nil,errors.Wrap(err,"Reset issue in userController")
	}
	return update,nil
}

func (m *Mongo)HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (m *Mongo)CheckPasswordHash(hash string,pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

