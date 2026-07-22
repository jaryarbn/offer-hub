package model

// PasswordRegisterReq represents the username/password registration request.
type PasswordRegisterReq struct {
	Username string `json:"username" form:"username" binding:"required,max=50"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

// PasswordRegisterResp is the unified response returned by password registration.
type PasswordRegisterResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// PasswordLoginReq represents the username/password login request.
type PasswordLoginReq struct {
	Username string `json:"username" form:"username" binding:"required,max=50"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

// PasswordLoginUserInfo is the userInfo object returned after login.
// Its JSON field names intentionally follow the camelCase authentication API contract.
type PasswordLoginUserInfo struct {
	UserID     string `json:"userId"`
	Username   string `json:"username"`
	NickName   string `json:"nickName"`
	Avatar     string `json:"avatar"`
	Sex        int    `json:"sex"`
	VIP        bool   `json:"vip"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	UserStatus int    `json:"userStatus"`
	UserType   int    `json:"userType"`
}

// PasswordLoginData is the data object returned after a successful login.
type PasswordLoginData struct {
	Token    string                `json:"token"`
	UserInfo PasswordLoginUserInfo `json:"userInfo"`
}

// PasswordLoginResp is the unified response returned by password login.
type PasswordLoginResp struct {
	Code int                `json:"code"`
	Msg  string             `json:"msg"`
	Data *PasswordLoginData `json:"data"`
}

// PasswordLogoutData contains the logout confirmation message.
type PasswordLogoutData struct {
	Message string `json:"message"`
}

// PasswordLogoutResp is the unified response returned by password logout.
type PasswordLogoutResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data *PasswordLogoutData `json:"data"`
}
