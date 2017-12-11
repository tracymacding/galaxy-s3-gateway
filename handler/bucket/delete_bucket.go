package bucket

import (
	"galaxy-s3-gateway/mux"
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/db"
	// "galaxy-s3-gateway/mongodb/dao"
	"galaxy-s3-gateway/gerror"
	"galaxy-s3-gateway/handler"
	"net/http"
)

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp handler.S3Responser

	defer func() {
		if err != nil {
			context.Set(r, "req_error", gerror.NewGError(err))
		}
		context.Set(r, "response", resp)
		// resp.Send(w)
	}()

	bucket := mux.Vars(r)["bucket"]

	err = db.ActiveService().DeleteBucket(bucket)
	if err != nil {
		if err == db.BucketNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchBucket", "/"+bucket)
		} else if err == db.BucketNotEmptyError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusConflict, r, "BucketNotEmpty", "/"+bucket)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/"+bucket)
		}
		return
	}
	resp = handler.NewS3NilResponse(http.StatusNoContent)
	return
}