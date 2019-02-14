package database

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"
	"stonesrv/conf"
	"stonesrv/env"
	"stonesrv/log"
	"stonesrv/models"
	"time"

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
	p.openCollection(env.CollectionUserBackup)
	p.openCollection(env.CollectionUpdates)
	p.openCollection(env.CollectionMAC)
	p.openCollection(env.CollectionDisk0)
	p.openCollection(env.CollectionUserBehavior)

	p.openCollection(env.CollectionPacks)

	p.openCollection(env.CollectionDailyInfo)

	p.openCollection(env.CollectionUserPack)
	p.openCollection(env.CollectionUserSpacePlus)
	p.openCollection(env.CollectionUserTablePlus)
	p.openCollection(env.CollectionUserSpaceAndTablePlus)

	p.openCollection(env.CollectionUserWeChat)
	p.openCollection(env.CollectionUserOrderWeChat)

	p.openCollection(env.CollectionOrders)
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
func (p *ArangoDB) RemoveMAC(userkey string) error {
	bindVar := map[string]interface{}{
		"userkey":     userkey,
		"@collection": env.CollectionMAC,
	}
	err := p.removeByUserKey(bindVar)
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
func (p *ArangoDB) RemoveDisk0(userkey string) error {
	bindVar := map[string]interface{}{
		"userkey":     userkey,
		"@collection": env.CollectionDisk0,
	}
	err := p.removeByUserKey(bindVar)
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
		log.Error(fmt.Sprintf("GetUserByKey() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Error(fmt.Sprintf("GetUserByKey() error, detail%v\n", err))
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
		"key":         user.Key,
		"@collection": env.CollectionUser,
	}
	err := p.removeByKey(bindVar)
	return err
}

//BackupUser 备份用户 带KEY
func (p *ArangoDB) BackupUser(user models.User) error {
	user.Key = ""
	bindVar := map[string]interface{}{
		"doc":         user,
		"@collection": env.CollectionUserBackup,
	}
	err := p.insert(bindVar)
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

//ExtendUser 延长用户时间 ext 单位是天
func (p *ArangoDB) ExtendUser(user models.User, ext int) error {
	loc, lok := time.LoadLocation("Local")
	regDate, pok := time.ParseInLocation(env.FullDateTimeFormat, user.RegDate, loc)
	if lok != nil || pok != nil {
		return fmt.Errorf("arangodb extenduser parse time failed")
	}
	expDate := regDate.Add(time.Hour * 24 * time.Duration(ext))
	user.ExpDate = expDate.Format(env.FullDateTimeFormat)
	user.Duration = int64(expDate.Sub(regDate).Minutes())
	return p.UpdateUserInfo(user)
}

//SetUserToPRO 用户升级为PRO用户
func (p *ArangoDB) SetUserToPRO(user models.User) error {
	user.Functions = 2
	return p.UpdateUserInfo(user)
}

//SetUserToSTD 用户降级为STD用户
func (p *ArangoDB) SetUserToSTD(user models.User) error {
	user.Functions = 1
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

//收费包
//PACK
//InsertPack
func (p *ArangoDB) InsertPack(userpack models.UserPack) error {
	bindVar := map[string]interface{}{
		"doc":         userpack,
		"@collection": env.CollectionUserPack,
	}
	err := p.insert(bindVar)
	return err
}

//GetPack
func (p *ArangoDB) GetPack(username string) *models.UserPack {
	user := p.GetUserByName(username)
	if user == nil {
		log.Error(fmt.Sprintf("GetPack() error, detail no user [%v] \n", username))
		return nil
	}

	cursor, err := p.queryDocByUserKey(env.CollectionUserPack, user.Key)
	if err != nil {
		log.Error(fmt.Sprintf("GetPack() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	var pack models.UserPack
	for cursor.HasMore() {
		_, err := cursor.ReadDocument(context.Background(), &pack)
		if err != nil {
			log.Error(fmt.Sprintf("GetPack() error, detail%v\n", err))
			return nil
		}
	}
	return &pack
}

//SPACE PLUS
//InsertUserSpacePlus
func (p *ArangoDB) InsertUserSpacePlus(userspaceplus models.UserSpacePlus) error {
	bindVar := map[string]interface{}{
		"doc":         userspaceplus,
		"@collection": env.CollectionUserSpacePlus,
	}
	err := p.insert(bindVar)
	return err
}

//GetUserSpacePlus
func (p *ArangoDB) GetUserSpacePlus(key string) []*models.UserSpacePlus {
	return nil
}

//TABLE PLUS
//InsertUserTablePlus
func (p *ArangoDB) InsertUserTablePlus(usertableplus models.UserTablePlus) error {
	bindVar := map[string]interface{}{
		"doc":         usertableplus,
		"@collection": env.CollectionUserTablePlus,
	}
	err := p.insert(bindVar)
	return err
}

//GetUserTablePlus
func (p *ArangoDB) GetUserTablePlus(key string) []*models.UserTablePlus {
	return nil
}

//SPACE AND TABLE
//InsertUserSpaceAndTablePlus
func (p *ArangoDB) InsertUserSpaceAndTablePlus(userspaceandtableplus models.UserSpaceAndTablePlus) error {
	bindVar := map[string]interface{}{
		"doc":         userspaceandtableplus,
		"@collection": env.CollectionUserSpaceAndTablePlus,
	}
	err := p.insert(bindVar)
	return err
}

//GetUserSpaceAndTablePlus
func (p *ArangoDB) GetUserSpaceAndTablePlus(key string) []*models.UserSpaceAndTablePlus {
	return nil
}

///////

//付款记录
//IsUserPaied 判断用户是否已经付款，现阶段都默认已经付款
func (p *ArangoDB) IsUserPaied(key string) bool {
	//return false
	return true
}

///用户行为记录
//UpsertUserBehavior 更新用户行为信息
func (p *ArangoDB) UpsertUserBehavior(userbehavior models.UserBehavior) error {
	bindVar := map[string]interface{}{
		"doc":         userbehavior,
		"key":         userbehavior.Key,
		"@collection": env.CollectionUserBehavior,
	}
	err := p.upsert(bindVar)
	return err
}

//GetUserBehaviorByKey 获取用户行为
func (p *ArangoDB) GetUserBehaviorByKey(key string) *models.UserBehavior {
	cursor, err := p.queryDocByKey(env.CollectionUserBehavior, key)
	if err != nil {
		log.Debug(fmt.Sprintf("GetUserBehaviorByKey() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var userbehavior models.UserBehavior
		_, err := cursor.ReadDocument(context.Background(), &userbehavior)
		if err != nil {
			log.Debug(fmt.Sprintf("GetUserBehaviorByKey() error, detail%v\n", err))
			return nil
		}
		return &userbehavior
	}
	return nil
}

//RecordUserPasswordFailed 记录用户登录密码错误次数
func (p *ArangoDB) RecordUserPasswordFailed(key string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.PasswordFailed++
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserInActivated 记录用户非激活状态登录次数
func (p *ArangoDB) RecordUserInActivated(key string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.InActivated++
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserExpired 记录用户过期登录次数
func (p *ArangoDB) RecordUserExpired(key string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.Expired++
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserLoginSuccess 记录用户最后登录成功时间和次数
func (p *ArangoDB) RecordUserLoginSuccess(key string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.LoginSuccess++
	userBehavior.LastLogin = time.Now().Format(env.FullDateTimeFormat)
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserLogoutSuccess 记录用户登出成功时间
func (p *ArangoDB) RecordUserLogoutSuccess(key string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.LastLogout = time.Now().Format(env.FullDateTimeFormat)
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserLoginIP 记录用户登录过的IP地址
func (p *ArangoDB) RecordUserLoginIP(key string, ip string) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.LoginIPs = append(userBehavior.LoginIPs, ip)
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserCurrentSpaces 记录用户已经使用的空间
func (p *ArangoDB) RecordUserCurrentSpaces(key string, space int64) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.UsedSpace = space
	return p.UpsertUserBehavior(*userBehavior)
}

//RecordUserCurrentTables 记录用户已经使用的表数量
func (p *ArangoDB) RecordUserCurrentTables(key string, table int64) error {
	userBehavior := p.GetUserBehaviorByKey(key)
	if userBehavior == nil {
		userBehavior = &models.UserBehavior{
			Key:      key,
			LoginIPs: make([]string, 0),
		}
	}
	userBehavior.UsedTable = table
	return p.UpsertUserBehavior(*userBehavior)
}

///用户行为记录

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
		"key":         update.Key,
		"@collection": env.CollectionUpdates,
	}
	err := p.removeByKey(bindVar)
	return err
}

//GetDailyInfo 获取每日一句
func (p *ArangoDB) GetDailyInfo() *models.DailyInfo {
	cursor, err := p.queryDailyInfoAll()
	if err != nil {
		log.Error(fmt.Sprintf("GetDailyInfo() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	dailyinfos := make([]*models.DailyInfo, 0)
	for cursor.HasMore() {
		var dailyInfo models.DailyInfo
		_, err := cursor.ReadDocument(context.Background(), &dailyInfo)
		if err != nil {
			log.Error(fmt.Sprintf("GetDailyInfo() error, detail%v\n", err))
			return nil
		}
		dailyinfos = append(dailyinfos, &dailyInfo)
	}
	//返回随机的每日一句
	alllen := len(dailyinfos)
	randIndex := rand.Intn(alllen)
	if randIndex >= 0 && randIndex < alllen {
		return dailyinfos[randIndex]
	}
	return nil
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
			SORT doc._key
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

//queryDocByUserKey 通用方法通过KEY找数据
func (p *ArangoDB) queryDocByUserKey(collection string, key string) (out godriver.Cursor, e error) {
	aql := `FOR doc in @@collection
			FILTER doc.userkey == @key
			SORT doc._key
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

//removeByKey 通用删除
func (p *ArangoDB) removeByKey(bindVar map[string]interface{}) error {
	aql := `
		REMOVE @key IN @@collection
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

//removeByUserKey 通用删除
func (p *ArangoDB) removeByUserKey(bindVar map[string]interface{}) error {
	aql := `
		FOR d IN @@collection
			FILTER d.userkey == @userkey
			REMOVE d IN @@collection
	`
	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		log.Info(fmt.Sprintf("REMOVE %v failed %v .", bindVar["@collection"], err))
		return err
	}
	defer cursor.Close()
	log.Info(fmt.Sprintf("REMOVE %v Success UserKey %v.", bindVar["@collection"], bindVar["userkey"]))
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

//queryDailyInfoAll 获取全部的每日一句
func (p *ArangoDB) queryDailyInfoAll() (out godriver.Cursor, e error) {
	aql := `FOR di in dailyinfo
				return di`

	cursor, err := p.Database.Query(context.Background(), aql, nil)
	if err != nil {
		return nil, err
	}
	return cursor, err
}
