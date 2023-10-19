package server

import (
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	cmds "github.com/ipfs/go-ipfs-cmds"

	"net/http"
)

type errResponse struct {
	Message string
	Code    cmds.ErrorType
	Type    string
}

func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	log.Infof("ErrorResponse %s: %v", r.URL.String(), err)
	v := errResponse{
		Code:    cmds.ErrNormal,
		Message: err.Error(),
		Type:    "error",
	}

	data, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err = w.Write(data); err != nil {
		log.Errorf("ErrorEncoder write response err:%v", err)
	}
}
