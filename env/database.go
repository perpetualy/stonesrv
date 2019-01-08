package env

var (
	//ArangoDBDefaultUser ArangoDB默认用户
	ArangoDBDefaultUser = "root"

	//ArangoDBDefaultPassword ArangoDB默认密码
	ArangoDBDefaultPassword = "827aOZ35vd"

	//ArangoDBDefaultDBName ArangoDB默认数据库名称
	ArangoDBDefaultDBName = "stone"

	//CollectionUser 用户集合
	CollectionUser = "user"

	//CollectionUpdates 更新集合
	CollectionUpdates = "updates"
	//CollectionMAC MAC地址集合
	CollectionMAC = "mac"
	//CollectionDisk0 Disk0地址集合
	CollectionDisk0 = "disk0"

	//CollectionUserBehavior 用户行为
	CollectionUserBehavior = "userbehavior"

	//CollectionPacks 套餐
	CollectionPacks = "packs"

	//CollectionUserPack 用户套餐
	CollectionUserPack = "userpack"
	//CollectionUserSpacePlus 空间叠加包
	CollectionUserSpacePlus = "userspaceplus"
	//CollectionUserTablePlus 表叠加包
	CollectionUserTablePlus = "usertableplus"
	//CollectionUserSpaceAndTablePlus 空间和表叠加包
	CollectionUserSpaceAndTablePlus = "userspaceandtableplus"

	//CollectionUserWeChat 用户和微信号关联
	CollectionUserWeChat = "userwechat"
	//CollectionUserOrderWeChat 用户订单和微信号
	CollectionUserOrderWeChat = "userorderwechat"

	//CollectionOrders 用户订单
	CollectionOrders = "orders"
)
