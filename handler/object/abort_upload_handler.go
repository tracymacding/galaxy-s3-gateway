package object

import (
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/db"
	"galaxy-s3-gateway/gerror"
	"galaxy-s3-gateway/handler"
	"galaxy-s3-gateway/mux"
	"errors"
	"net/http"
)

func AbortUploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp handler.S3Responser

	defer func() {
		if err != nil {
			context.Set(r, "req_error", gerror.NewGError(err))
		}
		context.Set(r, "response", resp)
		// resp.Send(w)
	}()

	if err = r.ParseForm(); err != nil {
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/")
		return
	}

	objectName := mux.Vars(r)["object"]
	bucket := mux.Vars(r)["bucket"]

	_, err = db.ActiveService().GetBucket(bucket)
	if err != nil {
		if err == db.BucketNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchBucket", "/"+bucket+objectName)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/"+bucket)
		}
		return
	}

	objectPath := bucket + "/" + objectName

	var uploadId string
	uploadId = r.Form.Get("uploadId")
	if uploadId == "" {
		err = errors.New("uploadId missing")
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusBadRequest, r, "InvalidParameter", "/")
		return
	}

	if err = validateUploadId(uploadId); err != nil {
		if err == db.UploadNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchUpload", "/" + objectPath)

		} else if err == AbortedUploadError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "UploadAborted", "/"+objectPath)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/"+objectPath)
		}
		return
	}

	if err = db.ActiveService().SetUploadAborted(uploadId); err != nil {
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", "/"+objectPath)
		return
	}
	resp = handler.NewS3NilResponse(http.StatusNoContent)
	return
}
