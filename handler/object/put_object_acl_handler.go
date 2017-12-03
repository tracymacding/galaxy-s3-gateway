package object

import (
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/handler"
	"galaxy-s3-gateway/gerror"
	"net/http"
)

// TODO: do nothing
func PutObjectACLHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp handler.S3Responser

	defer func() {
		if err != nil {
			context.Set(r, "req_error", gerror.NewGError(err))
		}
		// resp.Send(w)
		context.Set(r, "response", resp)
	}()

	resp = handler.NewS3NilResponse(http.StatusOK)
}
