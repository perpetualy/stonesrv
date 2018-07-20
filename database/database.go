package database

import "stonesrv/models"

type DataBase interface {
	Init()

	GetAllUsers() []*models.User
	GetUserByName(string) *models.User

	UpsertUser(models.User)
}

var db *DB

func Init() {
	db := &DB{}
	db.CurrentDB = &ArangoDB{}
	db.CurrentDB.Init()
}

type DB struct {
	CurrentDB DataBase
}

func GetDatabase() DataBase {
	return db.CurrentDB
}
