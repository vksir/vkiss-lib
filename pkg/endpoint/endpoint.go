package endpoint

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

type Handler interface {
	HttpMethod() string
	HttpRelativePath(relativePath string) string
	HttpHandler() gin.HandlerFunc
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// GetResource func indexResource(id string) (any, error)
type GetResource[P any] func(ctx context.Context, id string) (response P, err error)

func (f GetResource[P]) HttpMethod() string {
	return http.MethodGet
}

func (f GetResource[P]) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f GetResource[P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		resp, err := f(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

// ListResource func listResource(a any) ([]any, error)
type ListResource[Q, P any] func(ctx context.Context, request Q) (response []P, err error)

func (f ListResource[Q, P]) HttpMethod() string {
	return http.MethodGet
}

func (f ListResource[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f ListResource[Q, P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Q
		err := c.ShouldBindQuery(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		err = validate.Struct(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		resp, err := f(c, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

// CreateResource func createResource(a any) (any, error)
type CreateResource[Q, P any] func(ctx context.Context, request Q) (response P, err error)

func (f CreateResource[Q, P]) HttpMethod() string {
	return http.MethodPost
}

func (f CreateResource[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f CreateResource[Q, P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Q
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		err = validate.Struct(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		resp, err := f(c, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

// UpdateResource func updateResource(id string, a any) (any, error)
type UpdateResource[Q, P any] func(ctx context.Context, id string, request Q) (response P, err error)

func (f UpdateResource[Q, P]) HttpMethod() string {
	return http.MethodPut
}

func (f UpdateResource[Q, P]) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f UpdateResource[Q, P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req Q
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		err = validate.Struct(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		resp, err := f(c, id, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

// DeleteResource func deleteResource(id string) error
type DeleteResource func(ctx context.Context, id string) (err error)

func (f DeleteResource) HttpMethod() string {
	return http.MethodDelete
}

func (f DeleteResource) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f DeleteResource) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := f(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{})
	}
}

// GetOnlyOneResource func getResource() (any, error)
type GetOnlyOneResource[P any] func(ctx context.Context) (response P, err error)

func (f GetOnlyOneResource[P]) HttpMethod() string {
	return http.MethodGet
}

func (f GetOnlyOneResource[P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f GetOnlyOneResource[P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := f(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

// UpdateOnlyOneResource func updateResource(a any) (any, error)
type UpdateOnlyOneResource[Q, P any] func(ctx context.Context, request Q) (response P, err error)

func (f UpdateOnlyOneResource[Q, P]) HttpMethod() string {
	return http.MethodPut
}

func (f UpdateOnlyOneResource[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f UpdateOnlyOneResource[Q, P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Q
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		err = validate.Struct(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		resp, err := f(c, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, Response{Data: resp})
	}
}

type Resource struct {
	RelativePath string
	Handler      []Handler
}

func (r *Resource) LoadRouter(g *gin.RouterGroup) {
	for _, h := range r.Handler {
		g.Handle(h.HttpMethod(), h.HttpRelativePath(r.RelativePath), h.HttpHandler())
	}
}

type Api interface {
	Resource() []*Resource
}
