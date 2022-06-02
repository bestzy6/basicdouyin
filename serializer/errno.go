package serializer

type ErrNo int

const (
	OK             ErrNo = iota //正常
	ParamInvalid                // 参数不合法
	UserHasExisted              // 该 Username 已存在
	UserHasDeleted              // 用户已删除
	UserNotExisted              // 用户不存在
	WrongPassword               // 密码错误
	LoginRequired               // 用户未登录
	PermDenied                  // 没有操作权限

	// ...需要的话添加（by zy 2022年6月2日）

	UnknownError // 未知错误
)
