package kate

import (
	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
)

type ApiResponse struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"`
}

func Error(ctx context.Context, w ResponseWriter, errInfo ErrorInfo) {
	err := w.WriteJson(errInfo)
	if err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}

func Ok(ctx context.Context, w ResponseWriter) {
	OkData(ctx, w, nil)
}

func OkData(ctx context.Context, w ResponseWriter, data interface{}) {
	apiResp := &ApiResponse{
		ErrCode: 0,
		ErrMsg:  "success",
		Data:    data,
	}
	err := w.WriteJson(apiResp)
	if err != nil {
		log.Error(ctx, "write json response", "error", err)
	}
}
