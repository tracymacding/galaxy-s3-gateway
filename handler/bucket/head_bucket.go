package bucket

import (
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/mux"
	"galaxy-s3-gateway/db"
	// "galaxy-s3-gateway/mongodb/dao"
	"galaxy-s3-gateway/gerror"
	"galaxy-s3-gateway/handler"
	"net/http"
)

func HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp handler.S3Responser

	defer func() {
		if err != nil {
			context.Set(r, "req_error", gerror.NewGError(err))
		}
		// resp.Send(w)
		context.Set(r, "response", resp)
	}()

	bucket := mux.Vars(r)["bucket"]

	_, err = db.ActiveService().GetBucket(bucket)
	if err != nil {
		if err == db.BucketNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchBucket", "/"+bucket)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/"+bucket)
		}
		return
	}
	resp = handler.NewS3NilResponse(http.StatusOK)
	return
}
