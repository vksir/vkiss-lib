package ddns

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/netip"
)

func LoadRouter(g *gin.RouterGroup) {
	g.GET("/my_ip", getMyIp)
}

type MyIpResponse struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func getMyIp(r *gin.Context) {
	ap, err := netip.ParseAddrPort(r.Request.RemoteAddr)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		r.Abort()
		return
	}
	r.JSON(http.StatusOK, MyIpResponse{Ip: ap.Addr().String(), Port: int(ap.Port())})
}
