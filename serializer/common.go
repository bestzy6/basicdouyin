package serializer

// ErrNo 错误代码
type ErrNo int

const (
	OK             ErrNo = iota //正常
	ParamInvalid                // 参数不合法
	UserHasExisted              // 用户已存在
	UserHasDeleted              // 用户已删除
	UserNotExisted              // 用户不存在
	WrongPassword               // 密码错误
	LoginRequired               // 用户未登录
	PermDenied                  // 没有操作权限

	// ...需要其他错误码的话，在PermDenied下面添加（by zy 2022年6月2日）

	UnknownError ErrNo = 255 // 未知错误
)

type User struct {
	FollowCount   int64  `json:"follow_count"`   // 关注总数
	FollowerCount int64  `json:"follower_count"` // 粉丝总数
	ID            int64  `json:"id"`             // 用户id
	IsFollow      bool   `json:"is_follow"`      // true-已关注，false-未关注
	Name          string `json:"name"`           // 用户名称
}

type Comment struct {
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
	ID         int64  `json:"id"`          // 评论id
	User       User   `json:"user"`        // 评论用户信息
}

type Video struct {
	Author        User   `json:"author"`         // 视频作者信息
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	ID            int64  `json:"id"`             // 视频唯一标识
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`       // 视频播放地址
	Title         string `json:"title"`          // 视频标题
}
