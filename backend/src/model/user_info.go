package model

// UserInfo is the snake_case user information returned by the user-info API.
type UserInfo struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	NickName     string `json:"nick_name"`
	Avatar       string `json:"avatar"`
	VIP          bool   `json:"vip"`
	Sex          int    `json:"sex"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Introduction string `json:"introduction"`
	AvatarURL    string `json:"avatar_url"`
	UserStatus   int    `json:"user_status"`
	UserType     int    `json:"user_type"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}

// GetUserInfoResp is the unified response returned by the user-info API.
type GetUserInfoResp struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *UserInfo `json:"data"`
}
