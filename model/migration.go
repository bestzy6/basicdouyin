package model

//执行数据迁移

func migration() {
	// 自动迁移模式
	_ = DB.AutoMigrate(&User{})
	_ = DB.AutoMigrate(&Video{})
	_ = DB.AutoMigrate(&Post{})
	_ = DB.AutoMigrate(&FavoritePost{})
	_ = DB.AutoMigrate(&VideoCL{})
}
