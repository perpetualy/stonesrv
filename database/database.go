package database

import "stonesrv/models"

type DataBase interface {
	Init()

	GetAllUsers() []*models.User
	GetUserByName(string) *models.User
	
	UpsertUser(models.User)
}

var db = initDB()

func initDB() *DB{
	db := &DB{}
	db.CurrentDB = &ArangoDB{}
	db.CurrentDB.Init()
	return db
}

type DB struct{
	CurrentDB DataBase
}

func GetDatabase() DataBase{
	return db.CurrentDB
}