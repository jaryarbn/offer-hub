package comment

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"offer-hub/backend/src/model"
	backendservice "offer-hub/backend/src/service"
)

type commentServiceStub struct {
	addReq    model.AddCommentReq
	addUserID string
	addData   model.AddCommentData
	addErr    error
	addCalls  int
	deleteReq model.DeleteCommentReq
	deleteUID string
	deleteErr error
	deleteN   int
	updateReq model.UpdateCommentReq
	updateUID string
	updateOut model.UpdateCommentData
	updateErr error
	updateN   int
	req       model.ListCommentsReq
	userID    string
	data      model.ListCommentsData
	err       error
	calls     int
}

func (stub *commentServiceStub) DeleteComment(
	_ context.Context,
	req model.DeleteCommentReq,
	userID string,
) error {
	stub.deleteN++
	stub.deleteReq = req
	stub.deleteUID = userID
	return stub.deleteErr
}

func (stub *commentServiceStub) UpdateComment(
	_ context.Context,
	req model.UpdateCommentReq,
	userID string,
) (model.UpdateCommentData, error) {
	stub.updateN++
	stub.updateReq = req
	stub.updateUID = userID
	return stub.updateOut, stub.updateErr
}

func (stub *commentServiceStub) AddComment(
	_ context.Context,
	req model.AddCommentReq,
	userID string,
) (model.AddCommentData, error) {
	stub.addCalls++
	stub.addReq = req
	stub.addUserID = userID
	return stub.addData, stub.addErr
}

func (stub *commentServiceStub) ListComments(
	_ context.Context,
	req model.ListCommentsReq,
	userID string,
) (model.ListCommentsData, error) {
	stub.calls++
	stub.req = req
	stub.userID = userID
	return stub.data, stub.err
}

func TestAddCommentBindsBodyAndReturnsCompleteComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{addData: model.AddCommentData{
		CommentID: "comment-1",
		Comment: model.CommentInfo{
			CommentID:   "comment-1",
			UserID:      "user-1",
			UserName:    "Alice",
			UserAvatar:  "alice.png",
			Content:     "filtered content",
			ParentID:    "parent-1",
			ReplyTo:     "user-2",
			ReplyToName: "Bob",
			Status:      2,
			SubComments: []model.CommentInfo{},
		},
	}}
	controller := NewController(stub)
	engine := gin.New()
	engine.POST("/api/v1/comment/add", controller.AddComment)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/comment/add",
		bytes.NewBufferString(`{"target_type":1,"target_id":"question-1","parent_id":"parent-1","reply_to":"user-2","content":"original content"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("user_id", " user-1 ")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantReq := model.AddCommentReq{
		TargetType: 1,
		TargetID:   "question-1",
		ParentID:   "parent-1",
		ReplyTo:    "user-2",
		Content:    "original content",
	}
	if !reflect.DeepEqual(stub.addReq, wantReq) {
		t.Fatalf("bound request = %#v, want %#v", stub.addReq, wantReq)
	}
	if stub.addUserID != "user-1" {
		t.Fatalf("service user_id = %q, want user-1", stub.addUserID)
	}
	wantBody := `{"code":0,"msg":"success","data":{"comment_id":"comment-1","comment":{"comment_id":"comment-1","user_id":"user-1","user_name":"Alice","user_avatar":"alice.png","content":"filtered content","parent_id":"parent-1","reply_to":"user-2","reply_to_name":"Bob","status":2,"thumbs_up":0,"sub_comment_total":0,"user_liked":false,"sub_comments":[],"create_time":"","update_time":""}}}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestAddCommentRejectsUnauthenticatedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{}
	controller := NewController(stub)
	engine := gin.New()
	engine.POST("/api/v1/comment/add", controller.AddComment)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/comment/add",
		bytes.NewBufferString(`{"target_type":1,"target_id":"question-1","content":"content"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantBody := `{"code":401,"data":null,"msg":"未认证"}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
	if stub.addCalls != 0 {
		t.Fatalf("service calls = %d, want 0", stub.addCalls)
	}
}

func TestAddCommentRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		stubErr error
	}{
		{name: "unsupported target type", body: `{"target_type":3,"target_id":"question-1","content":"content"}`},
		{name: "missing content", body: `{"target_type":1,"target_id":"question-1"}`},
		{name: "blank content rejected by service", body: `{"target_type":1,"target_id":"question-1","content":"  "}`, stubErr: backendservice.ErrInvalidCommentContent},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			stub := &commentServiceStub{addErr: test.stubErr}
			controller := NewController(stub)
			engine := gin.New()
			engine.POST("/api/v1/comment/add", controller.AddComment)

			request := httptest.NewRequest(http.MethodPost, "/api/v1/comment/add", bytes.NewBufferString(test.body))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("user_id", "user-1")
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			wantBody := `{"code":400,"data":null,"msg":"invalid parameters"}`
			if response.Body.String() != wantBody {
				t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
			}
		})
	}
}

func TestDeleteCommentBindsBodyAndReturnsDocumentedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{}
	controller := NewController(stub)
	engine := gin.New()
	engine.POST("/api/v1/comment/delete", controller.DeleteComment)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/comment/delete",
		bytes.NewBufferString(`{"comment_id":"comment-1"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("user_id", " user-1 ")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if stub.deleteN != 1 || stub.deleteReq.CommentID != "comment-1" || stub.deleteUID != "user-1" {
		t.Fatalf("delete call = count %d, req %#v, user %q", stub.deleteN, stub.deleteReq, stub.deleteUID)
	}
	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	wantBody := `{"code":0,"msg":"success"}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestDeleteCommentMapsExpectedErrors(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		body     string
		stubErr  error
		wantBody string
		wantCall int
	}{
		{
			name:     "unauthenticated",
			body:     `{"comment_id":"comment-1"}`,
			wantBody: `{"code":401,"data":null,"msg":"未认证"}`,
		},
		{
			name:     "missing comment ID",
			userID:   "user-1",
			body:     `{}`,
			wantBody: `{"code":400,"data":null,"msg":"invalid parameters"}`,
		},
		{
			name:     "blank comment ID",
			userID:   "user-1",
			body:     `{"comment_id":" "}`,
			stubErr:  backendservice.ErrInvalidCommentID,
			wantBody: `{"code":400,"data":null,"msg":"invalid parameters"}`,
			wantCall: 1,
		},
		{
			name:     "another users comment",
			userID:   "user-1",
			body:     `{"comment_id":"comment-1"}`,
			stubErr:  backendservice.ErrCommentForbidden,
			wantBody: `{"code":403,"data":null,"msg":"forbidden"}`,
			wantCall: 1,
		},
		{
			name:     "comment not found",
			userID:   "user-1",
			body:     `{"comment_id":"comment-missing"}`,
			stubErr:  backendservice.ErrCommentNotFound,
			wantBody: `{"code":404,"data":null,"msg":"comment not found"}`,
			wantCall: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			stub := &commentServiceStub{deleteErr: test.stubErr}
			controller := NewController(stub)
			engine := gin.New()
			engine.POST("/api/v1/comment/delete", controller.DeleteComment)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/comment/delete",
				bytes.NewBufferString(test.body),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("user_id", test.userID)
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if response.Body.String() != test.wantBody {
				t.Fatalf("response body = %s, want %s", response.Body.String(), test.wantBody)
			}
			if stub.deleteN != test.wantCall {
				t.Fatalf("service calls = %d, want %d", stub.deleteN, test.wantCall)
			}
		})
	}
}

func TestUpdateCommentBindsBodyAndReturnsDocumentedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{updateOut: model.UpdateCommentData{CommentID: "comment-1"}}
	controller := NewController(stub)
	engine := gin.New()
	engine.POST("/api/v1/comment/update", controller.UpdateComment)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/comment/update",
		bytes.NewBufferString(`{"comment_id":"comment-1","content":"new content"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("user_id", " user-1 ")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantReq := model.UpdateCommentReq{CommentID: "comment-1", Content: "new content"}
	if stub.updateN != 1 || !reflect.DeepEqual(stub.updateReq, wantReq) || stub.updateUID != "user-1" {
		t.Fatalf("update call = count %d, req %#v, user %q", stub.updateN, stub.updateReq, stub.updateUID)
	}
	wantBody := `{"code":0,"msg":"success","data":{"comment_id":"comment-1"}}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestUpdateCommentMapsExpectedErrors(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		body     string
		stubErr  error
		wantBody string
		wantCall int
	}{
		{
			name:     "unauthenticated",
			body:     `{"comment_id":"comment-1","content":"new content"}`,
			wantBody: `{"code":401,"data":null,"msg":"未认证"}`,
		},
		{
			name:     "missing content",
			userID:   "user-1",
			body:     `{"comment_id":"comment-1"}`,
			wantBody: `{"code":400,"data":null,"msg":"invalid parameters"}`,
		},
		{
			name:     "blank content",
			userID:   "user-1",
			body:     `{"comment_id":"comment-1","content":" "}`,
			stubErr:  backendservice.ErrInvalidCommentContent,
			wantBody: `{"code":400,"data":null,"msg":"invalid parameters"}`,
			wantCall: 1,
		},
		{
			name:     "another users comment",
			userID:   "user-1",
			body:     `{"comment_id":"comment-1","content":"new content"}`,
			stubErr:  backendservice.ErrCommentForbidden,
			wantBody: `{"code":403,"data":null,"msg":"forbidden"}`,
			wantCall: 1,
		},
		{
			name:     "comment not found",
			userID:   "user-1",
			body:     `{"comment_id":"comment-missing","content":"new content"}`,
			stubErr:  backendservice.ErrCommentNotFound,
			wantBody: `{"code":404,"data":null,"msg":"comment not found"}`,
			wantCall: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			stub := &commentServiceStub{updateErr: test.stubErr}
			controller := NewController(stub)
			engine := gin.New()
			engine.POST("/api/v1/comment/update", controller.UpdateComment)

			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/comment/update",
				bytes.NewBufferString(test.body),
			)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("user_id", test.userID)
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if response.Body.String() != test.wantBody {
				t.Fatalf("response body = %s, want %s", response.Body.String(), test.wantBody)
			}
			if stub.updateN != test.wantCall {
				t.Fatalf("service calls = %d, want %d", stub.updateN, test.wantCall)
			}
		})
	}
}

func TestListCommentsBindsQueryAppliesDefaultsAndForwardsUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{data: model.ListCommentsData{
		List: make([]model.CommentInfo, 0),
	}}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/open/list_comments", controller.ListComments)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/open/list_comments?target_type=1&target_id=%20question-1%20",
		nil,
	)
	request.Header.Set("user_id", " user-1 ")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", response.Code, http.StatusOK)
	}
	wantReq := model.ListCommentsReq{
		TargetType:     1,
		TargetID:       "question-1",
		SortBy:         "create_time",
		SortOrder:      "desc",
		Page:           1,
		PageSize:       20,
		SubCommentPage: 1,
		SubCommentSize: 5,
	}
	if !reflect.DeepEqual(stub.req, wantReq) {
		t.Fatalf("bound request = %#v, want %#v", stub.req, wantReq)
	}
	if stub.userID != "user-1" {
		t.Fatalf("service user_id = %q, want user-1", stub.userID)
	}
	wantBody := `{"code":0,"msg":"success","data":{"total":0,"list":[]}}`
	if response.Body.String() != wantBody {
		t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
	}
}

func TestListCommentsBindsAllSupportedParameters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &commentServiceStub{data: model.ListCommentsData{List: []model.CommentInfo{}}}
	controller := NewController(stub)
	engine := gin.New()
	engine.GET("/api/v1/open/list_comments", controller.ListComments)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/open/list_comments?target_type=2&target_id=experience-1&parent_id=comment-1&sort_by=thumbs_up&sort_order=asc&page=2&page_size=10&sub_comment_page=3&sub_comment_size=4",
		nil,
	)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	wantReq := model.ListCommentsReq{
		TargetType:     2,
		TargetID:       "experience-1",
		ParentID:       "comment-1",
		SortBy:         "thumbs_up",
		SortOrder:      "asc",
		Page:           2,
		PageSize:       10,
		SubCommentPage: 3,
		SubCommentSize: 4,
	}
	if !reflect.DeepEqual(stub.req, wantReq) {
		t.Fatalf("bound request = %#v, want %#v", stub.req, wantReq)
	}
}

func TestListCommentsRejectsInvalidParameters(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{name: "missing target type", query: "target_id=question-1"},
		{name: "unsupported target type", query: "target_type=3&target_id=question-1"},
		{name: "missing target id", query: "target_type=1"},
		{name: "blank target id", query: "target_type=1&target_id=+++"},
		{name: "unsupported sort field", query: "target_type=1&target_id=question-1&sort_by=view_count"},
		{name: "unsupported sort direction", query: "target_type=1&target_id=question-1&sort_order=sideways"},
		{name: "negative page", query: "target_type=1&target_id=question-1&page=-1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			stub := &commentServiceStub{}
			controller := NewController(stub)
			engine := gin.New()
			engine.GET("/api/v1/open/list_comments", controller.ListComments)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/open/list_comments?"+test.query,
				nil,
			)
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			wantBody := `{"code":400,"data":null,"msg":"invalid parameters"}`
			if response.Body.String() != wantBody {
				t.Fatalf("response body = %s, want %s", response.Body.String(), wantBody)
			}
			if stub.calls != 0 {
				t.Fatalf("service calls = %d, want 0", stub.calls)
			}
		})
	}
}
