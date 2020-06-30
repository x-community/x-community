package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
)

const (
	// TracerContextKey represents opentracing tracer in context
	TracerContextKey = "X-Tracer-Context"
	// CurrentUserKey represents current login user in context
	CurrentUserKey = "X-Current-User"
)

// CurrentUser represents user info in context
type CurrentUser struct {
	ID   uint32
	User string
}

// SuccessResponse represents http success response body
type SuccessResponse struct{}

// ErrorResponse represents http error response
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents http api error
type APIError struct {
	Code    int32  `json:"code"`
	Message string `json:"message,omitempty"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Code: %d, Error: %s", e.Code, e.Message)
}

// APIService represents base type of api service
type APIService struct{}

// CurrentUser return current user info
func (*APIService) CurrentUser(ctx *gin.Context) *CurrentUser {
	return GetCurrentUser(ctx)
}

// UserID return current user id
func (s *APIService) UserID(ctx *gin.Context) (id uint32) {
	user := GetCurrentUser(ctx)
	if user != nil {
		id = user.ID
	}
	return id
}

// TracerContext return current opentracing span
func (s *APIService) TracerContext(ctx *gin.Context) context.Context {
	return TracerContext(ctx)
}

// Success represent success response
func (*APIService) Success(ctx *gin.Context, data ...interface{}) {
	if len(data) > 0 {
		ctx.JSON(http.StatusOK, data[0])
	} else {
		ctx.Status(http.StatusOK)
	}
}

// InvalidRequest represents invalid args error
func InvalidRequest(err string) APIError {
	return APIError{Code: 10001, Message: err}
}

// Failed represents failed http response
func (*APIService) Failed(ctx *gin.Context, code int, err ...APIError) {
	if len(err) > 0 {
		ctx.AbortWithStatusJSON(code, err[0])
	} else {
		ctx.AbortWithStatus(code)
	}
}

func (s *APIService) Error(ctx *gin.Context, err error) {
	if err == nil {
		s.Failed(ctx, http.StatusInternalServerError, ErrInternalServerError)
		return
	}
	var merr *merrors.Error
	if errors.As(err, &merr) && merr.Code != 500 {
		s.BadRequest(ctx, APIError{Code: merr.Code, Message: merr.Detail})
	} else {
		ctx.Error(err)
		s.Failed(ctx, http.StatusInternalServerError, ErrInternalServerError)
	}
}

// BadRequest represents bad request response
func (s *APIService) BadRequest(ctx *gin.Context, errs ...APIError) {
	for _, e := range errs {
		ctx.Error(e)
	}
	s.Failed(ctx, http.StatusBadRequest, errs...)
}

// Unauthorized represents unauthorized error
func (s *APIService) Unauthorized(ctx *gin.Context, err ...APIError) {
	for _, e := range err {
		ctx.Error(e)
	}
	s.Failed(ctx, http.StatusUnauthorized, err...)
}

// InternalServerError represents internal server error
func (s *APIService) InternalServerError(ctx *gin.Context, errs ...error) {
	for _, e := range errs {
		ctx.Error(e)
	}
	s.Failed(ctx, http.StatusInternalServerError, ErrInternalServerError)
}

// TracerContext return current opentracing span
func TracerContext(c *gin.Context) context.Context {
	v, exist := c.Get(TracerContextKey)
	if exist == false {
		return context.Background()
	}
	ctx, ok := v.(context.Context)
	if !ok {
		return context.Background()
	}
	return ctx
}

// GetTraceID return current tracing id
func GetTraceID(c *gin.Context) string {
	v, exist := c.Get(TracerContextKey)
	if exist {
		if ctx, ok := v.(context.Context); ok {
			md, ok := metadata.FromContext(ctx)
			if ok {
				return md["uber-trace-id"]
			}
		}
	}
	return ""
}

// GetCurrentUser return current user info
func GetCurrentUser(ctx *gin.Context) *CurrentUser {
	user, exists := ctx.Get(CurrentUserKey)
	if exists {
		currentUser, ok := user.(CurrentUser)
		if ok {
			return &currentUser
		}
	}
	return nil
}

// SetCurrentUser return set current user info
func SetCurrentUser(ctx *gin.Context, uid uint32, user string) {
	ctx.Set(CurrentUserKey, CurrentUser{ID: uid, User: user})
}

var (
	// ErrInternalServerError represents http 500 default error
	ErrInternalServerError = APIError{Code: 500, Message: "Internal Server Error"}
)
