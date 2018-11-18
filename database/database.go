package database

import "stonesrv/models"

//DataBase 数据库接口
type DataBase interface {
	Init()

	//硬件
	InsertMAC(models.MAC) error
	RemoveMAC(models.MAC) error
	InsertDisk0(models.Disk0) error
	RemoveDisk0(models.Disk0) error
	IsMACExist(mac string) bool
	IsDisk0Exist(disk0 string) bool

	//用户
	GetAllUsers() []*models.User
	IsUserExist(key string) bool
	GetUserByKey(key string) *models.User
	GetUserByName(username string) *models.User
	InsertUser(models.User) error
	RemoveUser(models.User) error
	//UpsertUser(models.User) error
	ActiveUser(models.User) error
	DeactiveUser(models.User) error
	ExtendUser(models.User, int) error
	UpdateUserInfo(models.User) error

	//用户行为
	UpsertUserBehavior(models.UserBehavior) error
	GetUserBehaviorByKey(key string) *models.UserBehavior
	RecordUserPasswordFailed(key string) error
	RecordUserInActivated(key string) error
	RecordUserExpired(key string) error
	RecordUserLoginSuccess(key string) error
	RecordUserLogoutSuccess(key string) error
	RecordUserLoginIP(key string, ip string) error
	RecordUserCurrentSpaces(key string, space int64) error
	RecordUserCurrentTables(key string, table int64) error

	//升级
	GetUpdate() *models.Updates
	SetUpdate(models.Updates) error
	RemoveUpdate(models.Updates) error
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
