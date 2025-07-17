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

type Endpoint interface {
	HttpMethod() string
	HttpRelativePath(relativePath string) string
	HttpHandler() gin.HandlerFunc
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Get func getResource() (any, error)
type Get[P any] func(ctx context.Context) (response P, err error)

func (f Get[P]) HttpMethod() string {
	return http.MethodGet
}

func (f Get[P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f Get[P]) HttpHandler() gin.HandlerFunc {
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

// GetByID func indexResource(id string) (any, error)
type GetByID[P any] func(ctx context.Context, id string) (response P, err error)

func (f GetByID[P]) HttpMethod() string {
	return http.MethodGet
}

func (f GetByID[P]) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f GetByID[P]) HttpHandler() gin.HandlerFunc {
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

// GetList func listResource(a any) ([]any, error)
type GetList[Q, P any] func(ctx context.Context, request Q) (response []P, err error)

func (f GetList[Q, P]) HttpMethod() string {
	return http.MethodGet
}

func (f GetList[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f GetList[Q, P]) HttpHandler() gin.HandlerFunc {
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

// Create func createResource(a any) (any, error)
type Create[Q, P any] func(ctx context.Context, request Q) (response P, err error)

func (f Create[Q, P]) HttpMethod() string {
	return http.MethodPost
}

func (f Create[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f Create[Q, P]) HttpHandler() gin.HandlerFunc {
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

// Update func updateResource(a any) (any, error)
type Update[Q, P any] func(ctx context.Context, request Q) (response P, err error)

func (f Update[Q, P]) HttpMethod() string {
	return http.MethodPut
}

func (f Update[Q, P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f Update[Q, P]) HttpHandler() gin.HandlerFunc {
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

// UpdateByID func updateResource(id string, a any) (any, error)
type UpdateByID[Q, P any] func(ctx context.Context, id string, request Q) (response P, err error)

func (f UpdateByID[Q, P]) HttpMethod() string {
	return http.MethodPut
}

func (f UpdateByID[Q, P]) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f UpdateByID[Q, P]) HttpHandler() gin.HandlerFunc {
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

// DeleteByID func deleteResource(id string) error
type DeleteByID func(ctx context.Context, id string) (err error)

func (f DeleteByID) HttpMethod() string {
	return http.MethodDelete
}

func (f DeleteByID) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f DeleteByID) HttpHandler() gin.HandlerFunc {
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

// DownloadByID func downloadFile(id string) error
type DownloadByID func(ctx context.Context, id string) (downloadPath string, afterDownload func(), err error)

func (f DownloadByID) HttpMethod() string {
	return http.MethodGet
}

func (f DownloadByID) HttpRelativePath(relativePath string) string {
	path, err := url.JoinPath(relativePath, ":id")
	if err != nil {
		panic(fmt.Sprintf("invalid relative path %s: %s", relativePath, err))
	}
	return path
}

func (f DownloadByID) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		downloadPath, afterDownload, err := f(c, id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		c.File(downloadPath)
		afterDownload()
	}
}

// Upload func uploadFile(filename string) (string, error)
type Upload[P any] func(ctx context.Context, filename string) (
	uploadPath string,
	afterUpload func() (response P, err error),
	err error)

func (f Upload[P]) HttpMethod() string {
	return http.MethodPost
}

func (f Upload[P]) HttpRelativePath(relativePath string) string {
	return relativePath
}

func (f Upload[P]) HttpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				Response{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		uploadPath, afterUpload, err := f(c, file.Filename)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		err = c.SaveUploadedFile(file, uploadPath)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				Response{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		resp, err := afterUpload()
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
	Handler      []Endpoint
}

func (r *Resource) LoadRouter(g *gin.RouterGroup) {
	for _, h := range r.Handler {
		g.Handle(h.HttpMethod(), h.HttpRelativePath(r.RelativePath), h.HttpHandler())
	}
}

type Api interface {
	Resource() []*Resource
}
