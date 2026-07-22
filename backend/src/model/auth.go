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

// ListCommentsReq represents the query parameters for listing comments.
type ListCommentsReq struct {
	TargetType     int    `json:"target_type" form:"target_type" binding:"required,oneof=1 2"`
	TargetID       string `json:"target_id" form:"target_id" binding:"required"`
	ParentID       string `json:"parent_id" form:"parent_id"`
	SortBy         string `json:"sort_by" form:"sort_by" binding:"omitempty,oneof=create_time thumbs_up"`
	SortOrder      string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
	Page           int    `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize       int    `json:"page_size" form:"page_size" binding:"omitempty,min=1"`
	SubCommentPage int    `json:"sub_comment_page" form:"sub_comment_page" binding:"omitempty,min=1"`
	SubCommentSize int    `json:"sub_comment_size" form:"sub_comment_size" binding:"omitempty,min=1"`
}

// CommentInfo is the comment object shared by comment list and mutation responses.
type CommentInfo struct {
	CommentID       string        `json:"comment_id"`
	UserID          string        `json:"user_id"`
	UserName        string        `json:"user_name"`
	UserAvatar      string        `json:"user_avatar"`
	Content         string        `json:"content"`
	ParentID        string        `json:"parent_id"`
	ReplyTo         string        `json:"reply_to"`
	ReplyToName     string        `json:"reply_to_name"`
	Status          int           `json:"status"`
	ThumbsUp        int           `json:"thumbs_up"`
	SubCommentTotal int64         `json:"sub_comment_total"`
	UserLiked       bool          `json:"user_liked"`
	SubComments     []CommentInfo `json:"sub_comments"`
	CreateTime      string        `json:"create_time"`
	UpdateTime      string        `json:"update_time"`
}

// ListCommentsData is the paginated data returned by the comment list endpoint.
type ListCommentsData struct {
	Total int64         `json:"total"`
	List  []CommentInfo `json:"list"`
}

// ListCommentsResp is the unified response returned by the comment list endpoint.
type ListCommentsResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data ListCommentsData `json:"data"`
}

// AddCommentReq represents a request to publish a top-level comment or reply.
type AddCommentReq struct {
	TargetType int    `json:"target_type" form:"target_type" binding:"required,oneof=1 2"`
	TargetID   string `json:"target_id" form:"target_id" binding:"required"`
	ParentID   string `json:"parent_id" form:"parent_id"`
	ReplyTo    string `json:"reply_to" form:"reply_to"`
	Content    string `json:"content" form:"content" binding:"required"`
}

// AddCommentData contains the newly created comment and its business ID.
type AddCommentData struct {
	CommentID string      `json:"comment_id"`
	Comment   CommentInfo `json:"comment"`
}

// AddCommentResp is the unified response returned after publishing a comment.
type AddCommentResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data AddCommentData `json:"data"`
}

// DeleteCommentReq represents a request to soft-delete a comment.
type DeleteCommentReq struct {
	CommentID string `json:"comment_id" form:"comment_id" binding:"required"`
}

// DeleteCommentResp follows the documented response, which has no data field.
type DeleteCommentResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// UpdateCommentReq represents a request to update the current user's comment.
type UpdateCommentReq struct {
	CommentID string `json:"comment_id" form:"comment_id" binding:"required"`
	Content   string `json:"content" form:"content" binding:"required"`
}

// UpdateCommentData identifies the comment updated by the request.
type UpdateCommentData struct {
	CommentID string `json:"comment_id"`
}

// UpdateCommentResp is the unified response returned after updating a comment.
type UpdateCommentResp struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data UpdateCommentData `json:"data"`
}

// GetCommentLikedCollectedCountData contains interaction totals for the current user's comments.
type GetCommentLikedCollectedCountData struct {
	LikedCount     int64 `json:"liked_count"`
	CollectedCount int64 `json:"collected_count"`
}

// GetCommentLikedCollectedCountResp is returned by the comment interaction-count endpoint.
type GetCommentLikedCollectedCountResp struct {
	Code int                               `json:"code"`
	Msg  string                            `json:"msg"`
	Data GetCommentLikedCollectedCountData `json:"data"`
}

// InteractionLikeReq represents a request to like a question or comment.
type InteractionLikeReq struct {
	TargetType int    `json:"target_type" form:"target_type" binding:"required,oneof=1 3"`
	TargetID   string `json:"target_id" form:"target_id" binding:"required"`
}

// InteractionLikeData contains the resulting like state and target like count.
type InteractionLikeData struct {
	Liked bool  `json:"liked"`
	Count int64 `json:"count"`
}

// InteractionLikeResp is the unified response returned after liking a target.
type InteractionLikeResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data InteractionLikeData `json:"data"`
}

// InteractionUnlikeReq has the same request contract as InteractionLikeReq.
type InteractionUnlikeReq = InteractionLikeReq

// InteractionUnlikeData contains the target like count after cancellation.
type InteractionUnlikeData struct {
	Count int64 `json:"count"`
}

// InteractionUnlikeResp is the unified response returned after cancelling a like.
type InteractionUnlikeResp struct {
	Code int                   `json:"code"`
	Msg  string                `json:"msg"`
	Data InteractionUnlikeData `json:"data"`
}

// TagQuestionReq represents a request to update the current user's question tag.
// Tag is a pointer so JSON binding can distinguish a missing field from the valid value 0.
type TagQuestionReq struct {
	QuestionID string `json:"question_id" form:"question_id" binding:"required"`
	Tag        *int   `json:"tag" form:"tag" binding:"required"`
}

// TagQuestionResp follows the documented response, whose data field is null.
type TagQuestionResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
