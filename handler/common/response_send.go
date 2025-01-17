package common

import (
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/handler"
	"net/http"
	"encoding/base64"
)

const (  
    base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"  
)

func base64Encode(src []byte) string {
	return base64.NewEncoding(base64Table).EncodeToString(src)
}

func SendResponseHandler(w http.ResponseWriter, r *http.Request) {
	reqId := context.Get(r, "req_id").(string)
	resp := context.Get(r, "response")
	if resp != nil {
		s3resp := resp.(handler.S3Responser)
		w.Header().Set("x-amz-id-2", reqId)
		w.Header().Set("x-amz-request-id", base64Encode([]byte(reqId)))
		w.Header().Set("Server", "GalaxyS3")
		s3resp.Send(w)
	}
}
