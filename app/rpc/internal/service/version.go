package service

import (
	"encoding/json"
	httpctx "github.com/go-kratos/kratos/v2/transport/http"
	"net/http"
)

type VersionService struct {
	info []byte
}

func NewVersionService(info *VersionInfo) *VersionService {
	b, _ := json.Marshal(info)
	return &VersionService{b}
}

func (s *VersionService) RegisterRoute(route *httpctx.Router) {
	route.POST("/version", s.Version)
}

type VersionInfo struct {
	Version string
	Commit  string
	Repo    string
	System  string
	Golang  string
}

func (s *VersionService) Version(ctx httpctx.Context) (err error) {
	w := ctx.Response()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write(s.info)
	return nil
}
