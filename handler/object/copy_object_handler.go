package object

import (
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/db"
	"galaxy-s3-gateway/fs"
	"galaxy-s3-gateway/gerror"
	"galaxy-s3-gateway/handler"
	"galaxy-s3-gateway/mux"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"io"
	"net/url"
	"encoding/xml"
	"fmt"
)

var (
	maxCopyObjectSize int64 = 5 * (1 << 30) 
)

type CopyObjectResult struct {
	XMLName       xml.Name `xml:"CopyObjectResult"`
	ETag          string   `xml:"ETag"`
	LastModified  string   `xml:"LastModified"`
}

func (resp *CopyObjectResult) Send(w http.ResponseWriter) { 
	body, _ := xml.MarshalIndent(resp, "", " ")
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func wrapS3CopyObjectResponse(lastModified int64, etag string) *CopyObjectResult {
	return &CopyObjectResult {
		ETag: etag,
		LastModified: time.Unix(lastModified, 0).Format(time.RFC3339),
	}
}

func parseSourceObject(r *http.Request) (string, string, error) {
	path := r.Header.Get("x-amz-copy-source")
	if path == "" {
		return "", "", errors.New("invalid x-amz-copy-source")
	}
	var err error
	// 根据AWS-S3定义,x-amz-copy-source为url encode后的
	path, err = url.QueryUnescape(path)
	if err != nil {
		return "", "", err
	}

	// x-amz-copy-source形式为"/bucket/object"
	entries := strings.SplitN(path, "/", 3)
	if len(entries) < 3 {
		return "", "", errors.New("invalid x-amz-copy-source")
	}
	return entries[1], entries[2], nil
}

func CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp handler.S3Responser

	defer func() {
		if err != nil {
			context.Set(r, "req_error", gerror.NewGError(err))
		}
		// resp.Send(w)
		context.Set(r, "response", resp)
	}()


	var srcBucket, srcObject string
	srcBucket, srcObject, err = parseSourceObject(r)
	if err != nil {
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusBadRequest, r, "InvalidParameter", "/")
		return
	}
	srcObjectPath := srcBucket + "/" + srcObject

	destBucket := mux.Vars(r)["bucket"]
	destObject := mux.Vars(r)["object"]
	destObjectFullPath := destBucket + "/" + destObject

	var destBucketObject *db.Bucket
	destBucketObject, err = db.ActiveService().GetBucket(destBucket)
	if err != nil {
		if err == db.BucketNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchBucket", destObjectFullPath)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", destObjectFullPath)
		}
		return
	}
	srcBucketObject := destBucketObject
	if srcBucket != destBucket {
		srcBucketObject, err = db.ActiveService().GetBucket(srcBucket)
		if err != nil {
			if err == db.BucketNotExistError {
				resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchBucket", srcObjectPath)
			} else {
				resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", srcObjectPath)
			}
			return
		}
	}

	var object *db.Object
	object, err = db.ActiveService().GetObject(srcBucket, srcObject, "", srcBucketObject.VersionEnabled)
	if err != nil {
		if err == db.ObjectNotExistError {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusNotFound, r, "NoSuchKey", srcObjectPath)
		} else {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", srcObjectPath)
		}
		return
	}

	if object.Size > maxCopyObjectSize {
		err = errors.New("can not copy object larger than 5GB")
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusBadRequest, r, "InvalidParameter", srcObjectPath)
		return
	}
	
	var reader io.Reader
	txId := context.Get(r, "req_id").(string)
	fid, etag := "", ""
	if object.Fid != "" {
		reader, err = fs.GetObject(object.Fid, txId)
		if err != nil {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", srcObjectPath)
			return
		}

		fid, etag, err = fs.PutObject(destBucketObject.UserID, object.Size, reader, txId)
		if err != nil {
			resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", srcObjectPath)
			return
		}
	}

	// keyMD5sum := md5.Sum(destObjectFullPath)
	objectMeta := &db.Object{
		ObjectName:    mux.Vars(r)["object"],
		Fid:           fid,
		//KeyMd5High:    md5.MD5High(keyMD5sum),
		//KeyMd5Low:     md5.MD5Low(keyMD5sum),
		//ConflictFlag:  0,
		//Md5High:       md5.MD5High([]byte(etag)),
		//Md5Low:        md5.MD5Low([]byte(etag)),
		Etag:          fmt.Sprintf("%s%s%s", "\"", etag, "\""),
		Bucket:        destBucket,
		Size:          object.Size,
		LastModified:  time.Now().Unix(),
		// DigestVersion: 0,
	}

	if err = db.ActiveService().PutObject(objectMeta, destBucketObject.VersionEnabled); err != nil {
		resp = handler.WrapS3ErrorResponseForRequest(http.StatusInternalServerError, r, "InternalError", srcObjectPath)
		return
	}
	resp = wrapS3CopyObjectResponse(objectMeta.LastModified, etag)
	return
}
