package main

import (
	"flag"
	"galaxy-s3-gateway/db"
	_ "galaxy-s3-gateway/db/mongodb"
	_ "galaxy-s3-gateway/db/mysql"
	"galaxy-s3-gateway/fs"
	"galaxy-s3-gateway/context"
	"galaxy-s3-gateway/handler/bucket"
	"galaxy-s3-gateway/handler/common"
	"galaxy-s3-gateway/handler/object"
	"galaxy-s3-gateway/middleware"
	"galaxy-s3-gateway/mux"
	"github.com/golang/glog"
	"net/http"
	_ "net/http/pprof"
)

func setupRouter() *mux.Router {

	router := mux.NewRouter()
	// generate req id before request, log request id after request
	// router.HookFunc(mux.HookBeforeRouter, common.GenerateRequestIdHandler).HookFunc(mux.HookAfterRouter, common.LogHandler)
	router.HookFunc(mux.HookAfterRouter, common.LogHandler).HookFunc(mux.HookAfterRouter, common.SendResponseHandler)

	// 支持AWS S3的virtual hosted–style request
	router.HandleFunc("/", bucket.ListBucketsHandler).Host("s3.galaxy.com").Methods("GET")
	// router.HandleFunc("/", bucket.GetBucketHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("GET")
	//router.HandleFunc("/", bucket.HeadBucketHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("HEAD")
	//router.HandleFunc("/", bucket.HeadBucketHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("HEAD")
	//router.HandleFunc("/", bucket.GetBucketACLHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("GET").Queries("acl", "{acl}")
	//router.HandleFunc("/", bucket.PutBucketHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("PUT")
	//router.HandleFunc("/", bucket.DeleteBucketHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("DELETE")
	//router.HandleFunc("/{object:.+}", object.InitMultipartUploadHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("POST").Queries("uploads", "{uploads}")
	//router.HandleFunc("/{object:.+}", object.ListPartsHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("GET").Queries("uploadId", "{uploadId}")
	//router.HandleFunc("/{object:.+}", object.UploadPartHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("PUT").Queries("uploadId", "{uploadId}", "partNumber", "{partNumber}")
	//router.HandleFunc("/{object:.+}", object.CompleteMultipartUploadHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("POST").Queries("uploadId", "{uploadId}")
	//router.HandleFunc("/{object:.+}", object.AbortUploadHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("DELETE").Queries("uploadId", "{uploadId}")
	//router.HandleFunc("/{object:.+}", object.CopyObjectHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").HeadersRegexp("x-amz-copy-source", ".+").Methods("PUT")
	//router.HandleFunc("/{object:.+}", object.PutObjectACLHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("PUT").Queries("acl", "{acl}")
	//router.HandleFunc("/{object:.+}", object.PutObjectHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("PUT")
	//router.HandleFunc("/{object:.+}", object.GetObjectACLHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("GET").Queries("acl", "{acl}")
	//router.HandleFunc("/{object:.+}", object.GetObjectHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("GET")
	//router.HandleFunc("/{object:.+}", object.HeadObjectHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("HEAD")
	//router.HandleFunc("/{object:.+}", object.DeleteObjectHandler).Host("{bucket:^[A-Za-z0-9_-]+$}.s3.galaxy.com").Methods("DELETE")
	router.HandleFunc("/", bucket.GetBucketHandler).Host("{bucket}.s3.galaxy.com").Methods("GET")
	router.HandleFunc("/", bucket.HeadBucketHandler).Host("{bucket}.s3.galaxy.com").Methods("HEAD")
	router.HandleFunc("/", bucket.HeadBucketHandler).Host("{bucket}.s3.galaxy.com").Methods("HEAD")
	router.HandleFunc("/", bucket.GetBucketACLHandler).Host("{bucket}.s3.galaxy.com").Methods("GET").Queries("acl", "{acl}")
	router.HandleFunc("/", bucket.PutBucketVersionHandler).Host("{bucket}.s3.galaxy.com").Methods("PUT").Queries("versioning", "{versioning}")
	router.HandleFunc("/", bucket.PutBucketHandler).Host("{bucket}.s3.galaxy.com").Methods("PUT")
	router.HandleFunc("/", bucket.DeleteBucketHandler).Host("{bucket}.s3.galaxy.com").Methods("DELETE")
	router.HandleFunc("/{object:.+}", object.InitMultipartUploadHandler).Host("{bucket}.s3.galaxy.com").Methods("POST").Queries("uploads", "{uploads}")
	router.HandleFunc("/{object:.+}", object.ListPartsHandler).Host("{bucket}.s3.galaxy.com").Methods("GET").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{object:.+}", object.UploadPartHandler).Host("{bucket}.s3.galaxy.com").Methods("PUT").Queries("uploadId", "{uploadId}", "partNumber", "{partNumber}")
	router.HandleFunc("/{object:.+}", object.CompleteMultipartUploadHandler).Host("{bucket}.s3.galaxy.com").Methods("POST").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{object:.+}", object.AbortUploadHandler).Host("{bucket}.s3.galaxy.com").Methods("DELETE").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{object:.+}", object.CopyObjectHandler).Host("{bucket}.s3.galaxy.com").HeadersRegexp("x-amz-copy-source", ".+").Methods("PUT")
	router.HandleFunc("/{object:.+}", object.PutObjectACLHandler).Host("{bucket}.s3.galaxy.com").Methods("PUT").Queries("acl", "{acl}")
	router.HandleFunc("/{object:.+}", object.PutObjectHandler).Host("{bucket}.s3.galaxy.com").Methods("PUT")
	router.HandleFunc("/{object:.+}", object.GetObjectACLHandler).Host("{bucket}.s3.galaxy.com").Methods("GET").Queries("acl", "{acl}")
	router.HandleFunc("/{object:.+}", object.GetObjectHandler).Host("{bucket}.s3.galaxy.com").Methods("GET").Queries("versionId", "{versionId}")
	router.HandleFunc("/{object:.+}", object.GetObjectHandler).Host("{bucket}.s3.galaxy.com").Methods("GET")
	router.HandleFunc("/{object:.+}", object.HeadObjectHandler).Host("{bucket}.s3.galaxy.com").Methods("HEAD")
	router.HandleFunc("/{object:.+}", object.DeleteObjectHandler).Host("{bucket}.s3.galaxy.com").Methods("DELETE")

	// 支持AWS S3的path-style request
	router.HandleFunc("/", bucket.ListBucketsHandler).Methods("GET")
	router.HandleFunc("/{bucket}/", bucket.HeadBucketHandler).Methods("HEAD")
	router.HandleFunc("/{bucket}/", bucket.GetBucketACLHandler).Methods("GET").Queries("acl", "{acl}")
	router.HandleFunc("/{bucket}/", bucket.GetBucketUploadsHandler).Methods("GET").Queries("uploads", "{uploads}")
	router.HandleFunc("/{bucket}/", bucket.GetBucketHandler).Methods("GET")
	router.HandleFunc("/{bucket}/", bucket.PutBucketVersionHandler).Methods("PUT").Queries("versioning", "{versioning}")
	router.HandleFunc("/{bucket}/", bucket.PutBucketHandler).Methods("PUT")
	router.HandleFunc("/{bucket}/", bucket.DeleteBucketHandler).Methods("DELETE")
	router.HandleFunc("/{bucket}", bucket.HeadBucketHandler).Methods("HEAD")
	router.HandleFunc("/{bucket}", bucket.GetBucketACLHandler).Methods("GET").Queries("acl", "{acl}")
	router.HandleFunc("/{bucket}", bucket.GetBucketHandler).Methods("GET")
	router.HandleFunc("/{bucket}", bucket.PutBucketVersionHandler).Methods("PUT").Queries("versioning", "{versioning}")
	router.HandleFunc("/{bucket}", bucket.PutBucketHandler).Methods("PUT")
	router.HandleFunc("/{bucket}", bucket.DeleteBucketHandler).Methods("DELETE")
	router.HandleFunc("/{bucket}/{object:.+}", object.InitMultipartUploadHandler).Methods("POST").Queries("uploads", "{uploads}")
	router.HandleFunc("/{bucket}/{object:.+}", object.ListPartsHandler).Methods("GET").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{bucket}/{object:.+}", object.UploadPartHandler).Methods("PUT").Queries("uploadId", "{uploadId}", "partNumber", "{partNumber}")
	router.HandleFunc("/{bucket}/{object:.+}", object.CompleteMultipartUploadHandler).Methods("POST").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{bucket}/{object:.+}", object.AbortUploadHandler).Methods("DELETE").Queries("uploadId", "{uploadId}")
	router.HandleFunc("/{bucket}/{object:.+}", object.CopyObjectHandler).HeadersRegexp("x-amz-copy-source", ".+").Methods("PUT")
	router.HandleFunc("/{bucket}/{object:.+}", object.PutObjectACLHandler).Methods("PUT").Queries("acl", "{acl}")
	router.HandleFunc("/{bucket}/{object:.+}", object.PutObjectHandler).Methods("PUT")
	router.HandleFunc("/{bucket}/{object:.+}", object.GetObjectACLHandler).Methods("GET").Queries("acl", "{acl}")
	router.HandleFunc("/{bucket}/{object:.+}", object.GetObjectHandler).Methods("GET").Queries("versionId", "{versionId}")
	router.HandleFunc("/{bucket}/{object:.+}", object.GetObjectHandler).Methods("GET")
	router.HandleFunc("/{bucket}/{object:.+}", object.HeadObjectHandler).Methods("HEAD")
	router.HandleFunc("/{bucket}/{object:.+}", object.DeleteObjectHandler).Methods("DELETE")
	return router
}

func startServe(router *mux.Router) {
	router.Use(&middleware.GenerateRequestIdMiddleware{})
	// router.Use(&middleware.AuthMiddleware{})
	// router.Use(&middleware.RateLimitMiddleware{})
	// http.ListenAndServe(":80", router)
	go func() {
		http.ListenAndServe(":3030", nil)
	}()
	http.ListenAndServe(*listenAddr, context.ClearHandler(router))
}

func main() {

	flag.Parse()

	// 使用mongodb作为后端存储,集群地址为192.168.100.100
	err := db.Open("mongodb", *mongoDBAddr)
	// 使用mysql作为后端存储,集群地址为192.168.100.100
	// err := db.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", *dbUser, *dbPassword, *dbAddr, *dbName))
	if err != nil {
		glog.Fatalf("init database error: %s", err)
	}

	err = fs.InitFS(*gfsMaster, *schedulerPath)
	if err != nil {
		glog.Fatalf("init file system error: %s", err)
	}

	startServe(setupRouter())
}
