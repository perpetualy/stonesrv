package database

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"stonesrv/conf"
	"stonesrv/env"
	"stonesrv/log"
	"stonesrv/models"

	godriver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"golang.org/x/net/context"
)

//ArangoDB 数据库
type ArangoDB struct {
	DataBase

	address  string
	user     string
	password string
	dbname   string
	err      error

	DbConn   godriver.Connection
	DbClient godriver.Client
	Database godriver.Database
}

//Init 初始化数据库
func (p *ArangoDB) Init() {
	//下面是硬编码
	p.address = fmt.Sprintf("http://%s:8529", conf.GetDBAddress())
	p.user = conf.GetDBUser()         //用户名
	p.password = conf.GetDBPassword() //密码
	p.dbname = conf.GetDBName()       //数据库名称

	p.initConnection(p.address)
	p.initClient(p.user, p.password)
	p.openDb(p.dbname)
	p.openCollection(env.CollectionUser)
	p.openCollection(env.CollectionUpdates)
	p.openCollection(env.CollectionMAC)
	p.openCollection(env.CollectionDisk0)
}

//InsertMAC 插入MAC
func (p *ArangoDB) InsertMAC(mac models.MAC) error {
	bindVar := map[string]interface{}{
		"doc":         mac,
		"@collection": env.CollectionMAC,
	}
	err := p.insert(bindVar)
	return err
}

//RemoveMAC 删除MAC
func (p *ArangoDB) RemoveMAC(mac models.MAC) error {
	bindVar := map[string]interface{}{
		"doc":         mac,
		"@collection": env.CollectionMAC,
	}
	err := p.remove(bindVar)
	return err
}

//InsertDisk0 插入Disk0
func (p *ArangoDB) InsertDisk0(disk0 models.Disk0) error {
	bindVar := map[string]interface{}{
		"doc":         disk0,
		"@collection": env.CollectionDisk0,
	}
	err := p.insert(bindVar)
	return err
}

//RemoveDisk0 删除Disk0
func (p *ArangoDB) RemoveDisk0(disk0 models.Disk0) error {
	bindVar := map[string]interface{}{
		"doc":         disk0,
		"@collection": env.CollectionDisk0,
	}
	err := p.remove(bindVar)
	return err
}

//IsMACExist 判断MAC是否存在
func (p *ArangoDB) IsMACExist(mac string) bool {
	cursor, err := p.queryDocByKey(env.CollectionMAC, mac)
	if err != nil {
		log.Error(fmt.Sprintf("IsMACExist() error, detail%v\n", err))
		return true
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var m models.MAC
		_, err := cursor.ReadDocument(context.Background(), &m)
		if err != nil {
			log.Debug(fmt.Sprintf("IsMACExist() error, detail%v\n", err))
			return true
		}
		return true
	}
	return false
}

//IsDisk0Exist 判断Disk0是否存在
func (p *ArangoDB) IsDisk0Exist(disk0 string) bool {
	cursor, err := p.queryDocByKey(env.CollectionDisk0, disk0)
	if err != nil {
		log.Error(fmt.Sprintf("IsDisk0Exist() error, detail%v\n", err))
		return true
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var d models.MAC
		_, err := cursor.ReadDocument(context.Background(), &d)
		if err != nil {
			log.Error(fmt.Sprintf("IsDisk0Exist() error, detail%v\n", err))
			return true
		}
		return true
	}
	return false
}

//GetAllUsers 获取所有的用户
func (p *ArangoDB) GetAllUsers() (users []*models.User) {
	cursor, err := p.queryUserAll()
	if err != nil {
		log.Error(fmt.Sprintf("GetAllUsers() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	users = make([]*models.User, 0)
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Error(fmt.Sprintf("GetAllUsers() error, detail%v\n", err))
			return nil
		}
		users = append(users, &user)
	}
	return users
}

//IsUserExist 判断用户是否存在
func (p *ArangoDB) IsUserExist(key string) bool {
	cursor, err := p.queryDocByKey(env.CollectionUser, key)
	if err != nil {
		log.Error(fmt.Sprintf("IsUserExist() error, detail%v\n", err))
		return true
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var u models.User
		_, err := cursor.ReadDocument(context.Background(), &u)
		if err != nil {
			log.Error(fmt.Sprintf("IsUserExist() error, detail%v\n", err))
			return true
		}
		return true
	}
	return false
}

//GetUserByKey 通过Key获取用户
func (p *ArangoDB) GetUserByKey(key string) *models.User {
	cursor, err := p.queryDocByKey(env.CollectionUser, key)
	if err != nil {
		log.Debug(fmt.Sprintf("GetUserByKey() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Debug(fmt.Sprintf("GetUserByKey() error, detail%v\n", err))
			return nil
		}
		return &user
	}
	return nil
}

//GetUserByName 通过名称获取用户
func (p *ArangoDB) GetUserByName(username string) *models.User {
	cursor, err := p.queryUserByName(username)
	if err != nil {
		log.Error(fmt.Sprintf("GetUserByName() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Error(fmt.Sprintf("GetUserByName() error, detail%v\n", err))
			return nil
		}
		return &user
	}
	return nil
}

//InsertUser 插入用户
func (p *ArangoDB) InsertUser(user models.User) error {
	bindVar := map[string]interface{}{
		"doc":         user,
		"@collection": env.CollectionUser,
	}
	err := p.insert(bindVar)
	return err
}

//RemoveUser 删除用户 带KEY
func (p *ArangoDB) RemoveUser(user models.User) error {
	bindVar := map[string]interface{}{
		"doc":         user,
		"@collection": env.CollectionUser,
	}
	err := p.remove(bindVar)
	return err
}

/* //UpsertUser 更新或者插入用户
func (p *ArangoDB) UpsertUser(user models.User) error {
	bindVar := map[string]interface{}{
		"doc":         user,
		"key":         user.Key,
		"@collection": env.CollectionUser,
	}
	err := p.upsert(bindVar)
	return err
} */

//ActiveUser 启用用户
func (p *ArangoDB) ActiveUser(user models.User) error {
	user.Activated = 1
	return p.UpdateUserInfo(user)
}

//DeactiveUser 关闭用户
func (p *ArangoDB) DeactiveUser(user models.User) error {
	user.Activated = 0
	return p.UpdateUserInfo(user)
}

//UpdateUserInfo 更新用户信息
func (p *ArangoDB) UpdateUserInfo(user models.User) error {
	bindVar := map[string]interface{}{
		"doc":         user,
		"key":         user.Key,
		"@collection": env.CollectionUser,
	}
	err := p.update(bindVar)
	return err
}

//GetUpdate 获取软件更新
func (p *ArangoDB) GetUpdate() *models.Updates {
	cursor, err := p.queryUpdate()
	if err != nil {
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var update models.Updates
		_, err := cursor.ReadDocument(context.Background(), &update)
		if err != nil {
			log.Error(fmt.Sprintf("GetUpdate() error, detail%v\n", err))
			return nil
		}
		return &update
	}
	return nil
}

//SetUpdate 设置更新
func (p *ArangoDB) SetUpdate(update models.Updates) error {
	bindVar := map[string]interface{}{
		"key":         update.Key,
		"doc":         update,
		"@collection": env.CollectionUpdates,
	}
	err := p.upsert(bindVar)
	return err
}

//RemoveUpdate 删除某一更新
func (p *ArangoDB) RemoveUpdate(update models.Updates) error {
	bindVar := map[string]interface{}{
		"doc":         update,
		"@collection": env.CollectionUpdates,
	}
	err := p.remove(bindVar)
	return err
}

//以下是内部函数

//initConnection 初始化连接
func (p *ArangoDB) initConnection(address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("arangodb.go : initConnection() %v %+v", err, string(debug.Stack())))
		}
	}()
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{address},
	})
	if err != nil {
		p.err = fmt.Errorf("arangodb create new connection error, detail%v", err)
	}
	p.DbConn = conn
	log.Info(fmt.Sprintf("DBConnection %v .", p.DbConn))
}

//initClient 初始化客户端
func (p *ArangoDB) initClient(user string, password string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("arangodb.go : initClient() %v %+v", err, string(debug.Stack())))
		}
	}()
	if p.err != nil {
		return
	}

	c, err := godriver.NewClient(godriver.ClientConfig{
		Connection:     p.DbConn,
		Authentication: godriver.BasicAuthentication(user, password),
	})
	if err != nil {
		p.err = fmt.Errorf("arangodb create new client error, detail%v", err)
	}
	p.DbClient = c
	log.Info(fmt.Sprintf("DBClient %v .", p.DbClient))
}

//openDb 打开数据库
func (p *ArangoDB) openDb(name string) {
	if p.err != nil {
		return
	}
	p.Database = p.ensureDatabase(nil, name, nil)
	log.Info(fmt.Sprintf("Database %v .", p.Database))
}

//openCollection 打开集合
func (p *ArangoDB) openCollection(name string) {
	if p.err != nil {
		return
	}
	p.ensureCollection(nil, name, nil)
	log.Info(fmt.Sprintf("Collection %v .", name))
}

//ensureDatabase 保证数据库存在
func (p *ArangoDB) ensureDatabase(ctx context.Context, name string, options *godriver.CreateDatabaseOptions) godriver.Database {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("arangodb.go : ensureDatabase() %v %+v", err, string(debug.Stack())))
		}
	}()
	if p.err != nil {
		return nil
	}

	db, err := p.DbClient.Database(ctx, name)
	if godriver.IsNotFound(err) {
		db, err = p.DbClient.CreateDatabase(ctx, name, options)
		if err != nil {
			if godriver.IsConflict(err) {
				p.err = fmt.Errorf("Failed to create database (conflict) '%s': %s %#v", name, p.describe(err), err)
			} else {
				p.err = fmt.Errorf("Failed to create database '%s': %s %#v", name, p.describe(err), err)
			}
		}
	} else if err != nil {
		p.err = fmt.Errorf("Failed to open database '%s': %s", name, p.describe(err))
	}
	return db
}

//ensureCollection 保证集合存在
func (p *ArangoDB) ensureCollection(ctx context.Context, name string, options *godriver.CreateCollectionOptions) godriver.Collection {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("arangodb.go : ensureCollection() %v %+v", err, string(debug.Stack())))
		}
	}()
	if p.err != nil {
		return nil
	}

	c, err := p.Database.Collection(ctx, name)
	if godriver.IsNotFound(err) {
		c, err = p.Database.CreateCollection(ctx, name, options)
		if err != nil {
			p.err = fmt.Errorf("Failed to create collection '%s': %s", name, p.describe(err))
		}
	} else if err != nil {
		p.err = fmt.Errorf("Failed to open collection '%s': %s", name, p.describe(err))
	}

	return c
}

// describe 返回错误字符串
func (p *ArangoDB) describe(err error) string {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("arangodb.go : describe() %v %+v", err, string(debug.Stack())))
		}
	}()
	if err == nil {
		return "nil"
	}
	cause := godriver.Cause(err)
	var msg string
	if re, ok := cause.(*godriver.ResponseError); ok {
		msg = re.Error()
	} else {
		c, _ := json.Marshal(cause)
		msg = string(c)
	}
	if cause.Error() != err.Error() {
		return fmt.Sprintf("%v caused by %v (%v)", err, cause, msg)
	}

	return fmt.Sprintf("%v (%v)", err, msg)

}

//queryUserConditions 获取用户的通用方法
func (p *ArangoDB) queryUserAll() (out godriver.Cursor, e error) {
	aql := `FOR u in user
				return u`

	cursor, err := p.Database.Query(context.Background(), aql, nil)
	if err != nil {
		return nil, err
	}
	return cursor, err
}

//queryUserByName 通过名称找用户数据
func (p *ArangoDB) queryUserByName(username string) (out godriver.Cursor, e error) {
	aql := `FOR u in user
			FILTER u.user == @UserName
			return u`
	bindVar := map[string]interface{}{
		"UserName": username,
	}
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		return nil, err
	}
	return cursor, err
}

//queryDocByKey 通用方法通过KEY找数据
func (p *ArangoDB) queryDocByKey(collection string, key string) (out godriver.Cursor, e error) {
	aql := `FOR doc in @@collection
			FILTER doc._key == @key
			return doc`
	bindVar := map[string]interface{}{
		"key":         key,
		"@collection": collection,
	}
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		return nil, err
	}
	return cursor, err
}

//insert 通用插入
func (p *ArangoDB) insert(bindVar map[string]interface{}) error {
	aql := `
		INSERT @doc IN @@collection
		RETURN NEW._key
		`
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		log.Info(fmt.Sprintf("INSERT %v failed %v .", bindVar["@collection"], err))
		return err
	}
	defer cursor.Close()
	var ret string
	_, err = cursor.ReadDocument(context.Background(), &ret)
	if err != nil {
		log.Info(fmt.Sprintf("INSERT %v failed %v .", bindVar["@collection"], err))
		return err
	}
	log.Info(fmt.Sprintf("INSERT %v Success %v.", bindVar["@collection"], ret))
	return nil
}

//remove 通用删除
func (p *ArangoDB) remove(bindVar map[string]interface{}) error {
	aql := `
		REMOVE @doc IN @@collection
		RETURN OLD._key
	`
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		log.Info(fmt.Sprintf("REMOVE %v failed %v .", bindVar["@collection"], err))
		return err
	}
	defer cursor.Close()
	var ret string
	_, err = cursor.ReadDocument(context.Background(), &ret)
	if err != nil {
		log.Info(fmt.Sprintf("REMOVE %v failed %v .", bindVar["@collection"], err))
		return err
	}
	log.Info(fmt.Sprintf("REMOVE %v Success %v.", bindVar["@collection"], ret))
	return nil
}

//update 通用更新
func (p *ArangoDB) update(bindVar map[string]interface{}) error {
	aql := `
			UPDATE @key WITH @doc IN @@collection
			RETURN NEW._key
	`
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		log.Info(fmt.Sprintf("UPDATE %v failed %v .", bindVar["@collection"], err))
		return err
	}
	defer cursor.Close()
	var ret string
	_, err = cursor.ReadDocument(context.Background(), &ret)
	if err != nil {
		log.Info(fmt.Sprintf("UPDATE %v failed %v .", bindVar["@collection"], err))
		return err
	}
	log.Info(fmt.Sprintf("UPDATE %v Success %v.", bindVar["@collection"], ret))
	return nil
}

//upsert 通用插入与更新
func (p *ArangoDB) upsert(bindVar map[string]interface{}) error {
	aql := `
			UPSERT { _key: @key }
				INSERT @doc
				UPDATE @doc IN @@collection
			LET opType = IS_NULL(OLD) ? "insert" : "update"
			RETURN { type: opType }
	`
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		log.Info(fmt.Sprintf("UPSERT %v failed %v .", bindVar["@collection"], err))
		return err
	}
	defer cursor.Close()
	ret := struct {
		Type string
	}{}
	_, err = cursor.ReadDocument(context.Background(), &ret)
	if err != nil {
		log.Info(fmt.Sprintf("UPSERT %v failed %v .", bindVar["@collection"], err))
		return err
	}
	log.Info(fmt.Sprintf("UPSERT %v Success %v Type %v.", bindVar["@collection"], bindVar["key"], ret.Type))
	return nil
}

//查询排序过后的update
func (p *ArangoDB) queryUpdate() (cursor godriver.Cursor, err error) {
	aql := `FOR upd in updates
			SORT upd.reldate
			return upd`
	cursor, err = p.Database.Query(context.Background(), aql, nil)
	return
}
