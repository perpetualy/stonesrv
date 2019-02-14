package accounts

import "sync"

var acc = newAccount()

//newAccount 新建本地认证账户对象
func newAccount() *Accounts {
	accounts := &Accounts{
		account: make(map[string]string),
	}

	//添加系统默认账户
	accounts.addAccount("xmvideo", "456789") //xmvideo admin 账户

	return accounts
}

//Accounts 本地认证账户
type Accounts struct {
	mutex   sync.Mutex
	account map[string]string
}

//AddAccount 添加本地认证账户
func AddAccount(user string, password string) string {
	acc.addAccount(user, password)
	return user
}

//GetAccounts 获取全部认证账户
func GetAccounts() map[string]string {
	return acc.account
}

//addAccount 添加单个账户
func (p *Accounts) addAccount(user string, password string) {
	p.mutex.Lock()
	p.account[user] = password
	p.mutex.Unlock()
}
