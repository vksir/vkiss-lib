package endpoint

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
)

type api struct{}

func (a *api) Resource() []Resource {
	return []Resource{
		{
			RelativePath: "/api/resource",
			Handler: []Endpoint{
				GetByID[resourceResp](a.getResource),
				GetList[listResourceReq, resourceResp](a.listResource),
				Create[createResourceReq, resourceResp](a.createResource),
				UpdateByID[updateResourceReq, resourceResp](a.updateResource),
				DeleteByID(a.deleteResource),
			},
		},
	}
}

type resourceResp struct {
	Id string `json:"id"`
}
type listResourceReq struct {
	Name string `json:"name"`
}
type createResourceReq struct {
	Name string `json:"name" validate:"required"`
}
type updateResourceReq struct {
	Name string `json:"name"`
}

func (a *api) getResource(ctx context.Context, id string) (resourceResp, error) {
	fmt.Println("get")
	return resourceResp{Id: "1"}, nil
}

func (a *api) listResource(ctx context.Context, req listResourceReq) ([]resourceResp, error) {
	fmt.Println("list:", req)
	return []resourceResp{{Id: "1"}, {Id: "2"}}, nil
}

func (a *api) createResource(ctx context.Context, req createResourceReq) (resourceResp, error) {
	fmt.Println("create:", req)
	return resourceResp{"1"}, nil
}

func (a *api) updateResource(ctx context.Context, id string, req updateResourceReq) (resourceResp, error) {
	fmt.Println("update:", req)
	return resourceResp{id}, nil
}

func (a *api) deleteResource(ctx context.Context, id string) error {
	fmt.Println("delete:", id)
	return nil
}

func TestResource(t *testing.T) {
	e := gin.New()
	a := &api{}
	for _, r := range a.Resource() {
		r.LoadRouter(&e.RouterGroup)
	}
	fmt.Println(e.Run(":5800"))
}
