package database

import "stonesrv/models"

//DataBase 数据库接口
type DataBase interface {
	Init()

	GetAllUsers() []*models.User
	GetUserByKey(string) *models.User
	GetUserByName(string) *models.User

	UpsertUser(models.User)
}

var db *DB

//Init 初始化数据库
func Init() {
	db = &DB{}
	db.CurrentDB = &ArangoDB{}
	db.CurrentDB.Init()
}

//DB 数据库结构
type DB struct {
	CurrentDB DataBase
}

//GetDatabase 获取当前数据库
func GetDatabase() DataBase {
	return db.CurrentDB
}
