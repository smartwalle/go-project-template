package pkg

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/nhttp"
	"io"
	"os"
	"path"
)

type AliOSSClient struct {
	conf   AliOSSConfig
	client *oss.Client
}

func NewAliOSSClient(conf AliOSSConfig) *AliOSSClient {
	var nClient = &AliOSSClient{}
	nClient.conf = conf
	client, err := oss.New(conf.Endpoint, conf.Key, conf.Secret)
	if err != nil {
		log4go.Errorln("初始化 ali OSS 发生错误:", err)
		os.Exit(-1)
	}
	nClient.client = client
	return nClient
}

func (this *AliOSSClient) Upload(file io.Reader, name string) (result string, err error) {
	bucket, err := this.client.Bucket(this.conf.Bucket)
	if err != nil {
		return "", err
	}

	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	objectAcl := oss.ObjectACL(oss.ACLDefault)

	fileExt := path.Ext(name)

	uuidName := uuid.New().String() + fileExt

	nPath := path.Join(this.conf.Bucket, uuidName)

	// 上传数据
	err = bucket.PutObject(nPath, file, storageType, objectAcl)
	if err != nil {
		return "", err
	}

	var nURL = nhttp.MustURL(this.conf.Domain)
	nURL.JoinPath(nPath)

	return nURL.String(), nil
}

func (this *AliOSSClient) UploadFile(file string, name string) (string, error) {
	bucket, err := this.client.Bucket(this.conf.Bucket)
	if err != nil {
		return "", err
	}

	// 指定存储类型为归档存储。
	storageType := oss.ObjectStorageClass(oss.StorageStandard)

	// 指定访问权限为公共读。
	objectAcl := oss.ObjectACL(oss.ACLDefault)

	var nName = ""
	if name == "" {
		fileExt := path.Ext(file)
		nName = uuid.New().String() + fileExt
	} else {
		_, nName = path.Split(file)
	}

	nPath := path.Join(this.conf.Bucket, nName)

	// 上传数据
	err = bucket.PutObjectFromFile(nPath, file, storageType, objectAcl)
	if err != nil {
		return "", err
	}

	var nURL = nhttp.MustURL(this.conf.Domain)
	nURL.JoinPath(nPath)

	return nURL.String(), nil
}
