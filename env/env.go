package env

var (
	//FullDateTimeFormat 时间格式字符串
	FullDateTimeFormat = "2006-01-02 15:04:05"
)

//回报消息类型定义
var (
	ParamsErrors = 400
	AuthFailed   = 401
	PageNotFound = 404

	RegSuccess = 600
	//	RegFailedParamsErrors          = 601
	RegFailedUserAlreadyRegistered = 602
	RegFailedPCAlreadyRegistered   = 603
	RegFailedInvalidDuration       = 604
	RegFailed                      = 605

	LoginSuccess = 610
	//	LoginFailedParamsErrors      = 611
	LoginFailedUserDoNotExists   = 612
	LoginFailedPasswordIncorrect = 613
	LoginFailedUserInactivated   = 614
	LoginFailedGetDateFailed     = 615
	LoginFailedUserExpired       = 616
	LoginFailedGenTokenFailed    = 617

	LogoutSuccess = 620
	//	LogoutFailedParamsErrors    = 621
	LogoutFailedUserDoNotExists = 622
	LogoutFailed                = 623

	GetUserInfoSuccess = 630
	//	GetUserInfoFailedParamsErrors    = 631
	GetUserInfoFailedUserDoNotExists = 632
	GetUserInfoFailed                = 633

	GetUpdatesUpdateFound = 640
	GetUpdatesNoNeed      = 641
	GetUpdatesEmergent    = 642
	//	GetUpdatesFailedParamsError    = 643
	GetUpdatesFailedCheckingFailed = 644
	GetUpdatesFailedRemoteFailed   = 645
	GetUpdatesLocalUpdateAlready   = 646

	GetClassifiesSuccess   = 650
	GetClassifyInfoSuccess = 651

	AddTableSuccess             = 660
	AddTableFailedAlreadyExists = 661

	UpdateTableSuccess           = 670
	UpdateTableFailedDoNotExists = 671

	GetIndexsRecommendsSuccess           = 680
	GetIndexsRecommendsFailedDataInvalid = 681
	GetIndexsRecommendsFailedNotFound    = 682

	InsertTableSuccess = 690

	GetAllTablesSuccess = 700
)
