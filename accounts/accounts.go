package accounts

import "sync"

var acc = newAccount()

func newAccount() *Accounts{
	accounts := &Accounts{
		account :make(map[string]string),
	}

	//添加系统默认账户
	accounts.addAccount("stone", "456789")//stone admin 账户

	return accounts
}

type Accounts struct {
	mutex sync.Mutex
	account map[string]string
}

func AddAccount(user string, password string){
	acc.addAccount(user, password)
}

func GetAccounts() map[string]string{
	return acc.account
}

func (p *Accounts) addAccount(user string, password string){
	p.mutex.Lock()
	p.account[user] = password
	p.mutex.Unlock()
}