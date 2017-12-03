package fs

import (
	"fmt"
	sdk "galaxyfs-go-sdk"
	"io"
	"strconv"
	"strings"
	"time"
)

var fs sdk.FS

func InitFS(master string) error {
	var err error
	opt := &sdk.Option{
		ZKAddr:    master,
		Timeout:      time.Second * 10,
		ConnPoolSize: 10,
	}

	fs, err = sdk.NewGFS(opt)
	if err != nil {
		return err
	}
	return nil
	// fs = &sdk.DummyFS{}
	// return nil
}

// PutObject write object content to underlying Galaxy distributed file system
func PutObject(user string, size int64, input io.Reader, reqId string) (string, string, error) {
	putRequest := &sdk.PutRequest{
		User:     user,
		Disk:     "sata",
		Replicas: 3,
		Size:     size,
		Input:    input,
		TxnId:    reqId,
	}
	resp := fs.Put(putRequest)
	if resp.Error != nil {
		return "", "", resp.Error
	}
	return resp.Fid.String(), resp.MD5, nil
}

// DeleteObject delete object from underlying Galaxy distributed file system
func DeleteObject(fid string, reqId string) error {
	ids := strings.Split(fid, "-")
	if len(ids) != 2 {
		return fmt.Errorf("invalid file id %s", fid)
	}

	blockId, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file id %s", fid)
	}
	internalId, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file id %s", fid)
	}

	delRequest := &sdk.DeleteRequest{
		Fid:   sdk.GFileId{BlockId: blockId, InternalId: internalId},
		TxnId: reqId,
	}
	resp := fs.Delete(delRequest)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

// GetObject get object content from underlying galaxy distributed file system
func GetObject(fid string, reqId string) (io.ReadSeeker, error) {
	ids := strings.Split(fid, "-")
	if len(ids) != 2 {
		return nil, fmt.Errorf("invalid file id %s", fid)
	}

	blockId, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file id %s", fid)
	}
	internalId, err := strconv.ParseInt(ids[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file id %s", fid)
	}

	getRequest := &sdk.GetRequest{
		Fid:   sdk.GFileId{BlockId: blockId, InternalId: internalId},
		TxnId: reqId,
	}
	resp := fs.Get(getRequest)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Output, nil
}
