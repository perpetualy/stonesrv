package database

import "stonesrv/models"

//DataBase 数据库接口
type DataBase interface {
	Init()

	//硬件
	InsertMAC(models.MAC) error
	InsertDisk0(models.Disk0) error
	IsMACExist(mac string) bool
	IsDisk0Exist(disk0 string) bool

	//用户
	GetAllUsers() []*models.User
	IsUserExist(key string) bool
	GetUserByKey(string) *models.User
	GetUserByName(string) *models.User
	InsertUser(models.User) error
	UpsertUser(models.User) error

	//升级
	GetUpdate() *models.Update
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
