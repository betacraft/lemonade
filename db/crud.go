package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func MgCreate(c string, data interface{}) error {
	session := GetMongoSession()
	defer session.Close()
	return session.DB(dbName).C(c).Insert(data)
}

func MgCreateStrong(c string, data interface{}) error {
	session := GetMongoSession()
	defer session.Close()
	session.SetMode(mgo.Strong, true)
	return session.DB(dbName).C(c).Insert(data)
}

func MgUpdate(c string, id bson.ObjectId, data interface{}) error {
	session := GetMongoSession()
	defer session.Close()
	return session.DB(dbName).C(c).Update(bson.M{"_id": id}, data)
}

func MgUpdateStrong(c string, id bson.ObjectId, data interface{}) error {
	session := GetMongoSession()
	session.SetMode(mgo.Strong, true)
	defer session.Close()
	return session.DB(dbName).C(c).Update(bson.M{"_id": id}, data)
}

func MgRetrieve(c string, data interface{}) error {
	session := GetMongoSession()
	defer session.Close()
	return session.DB(dbName).C(c).Find(nil).All(data)
}

func MgRetrieveStrong(c string, data interface{}) error {
	session := GetMongoSession()
	session.SetMode(mgo.Strong, true)
	defer session.Close()
	return session.DB(dbName).C(c).Find(nil).All(data)
}

func MgFindOne(c string, find *bson.M, data interface{}) error {
	session := GetMongoSession()
	defer session.Close()
	return session.DB(dbName).C(c).Find(find).One(data)
}

func MgFindOneStrong(c string, find *bson.M, data interface{}) error {
	session := GetMongoSession()
	session.SetMode(mgo.Strong, true)
	defer session.Close()
	return session.DB(dbName).C(c).Find(find).One(data)
}

func MgFindAll(c string, find *bson.M, data interface{}) error {
	db := GetMongo()
	defer db.Session.Close()
	return db.C(c).Find(find).All(data)
}

func MgFindPage(c string, find *bson.M, page int, data interface{}) error {
	db := GetMongo()
	defer db.Session.Close()
	return db.C(c).Find(find).Skip(page * 9).Limit(9).All(data)
}
