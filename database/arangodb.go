package database

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"stonesrv/conf"
	"stonesrv/env"
	"stonesrv/log"
	"stonesrv/models"
	"strings"

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
	p.user = "root"           //用户名
	p.password = "827aOZ35vd" //密码
	p.dbname = "stone"        //数据库名称

	p.initConnection(p.address)
	p.initClient(p.user, p.password)
	p.openDb(p.dbname)
	p.openCollection(env.CollectionUser)
	p.openCollection(env.CollectionUpdate)
}

//GetAllUsers 获取所有的用户
func (p *ArangoDB) GetAllUsers() (users []*models.User) {
	cursor, err := p.queryUserConditions("")
	if err != nil {
		log.Debug(fmt.Sprintf("GetAllUsers() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	users = make([]*models.User, 0)
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Debug(fmt.Sprintf("GetAllUsers() error, detail%v\n", err))
			return nil
		}
		users = append(users, &user)
	}
	return users
}

//GetUserByName 通过名称获取用户
func (p *ArangoDB) GetUserByName(username string) *models.User {
	cursor, err := p.queryUserConditions(username)
	if err != nil {
		log.Debug(fmt.Sprintf("GetUserByName() error, detail%v\n", err))
		return nil
	}
	defer cursor.Close()
	for cursor.HasMore() {
		var user models.User
		_, err := cursor.ReadDocument(context.Background(), &user)
		if err != nil {
			log.Debug(fmt.Sprintf("GetUserByName() error, detail%v\n", err))
			return nil
		}
		return &user
	}
	return nil
}

//UpsertUser 更新或者插入用户
func (p *ArangoDB) UpsertUser(user models.User) {
	bindVar := map[string]interface{}{
		"doc":         user,
		"key":         user.Key,
		"@collection": env.CollectionUser,
	}
	err := p.upsert(bindVar)
	if err == nil {
	}
}

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
func (p *ArangoDB) queryUserConditions(username string) (out godriver.Cursor, e error) {
	aql := `FOR u in user
				return u`
	bindVar := map[string]interface{}{}
	if strings.Compare(strings.Trim(username, " "), "") != 0 {
		aql = `FOR u in user
				FILTER u.user == @UserName
				return u`
		bindVar = map[string]interface{}{
			"UserName": username,
		}
	}

	cursor, err := p.Database.Query(context.Background(), aql, bindVar)
	if err != nil {
		return nil, err
	}
	return cursor, err
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
