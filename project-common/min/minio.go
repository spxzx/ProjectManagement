package min

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strconv"
)

type MinioClient struct {
	c *minio.Client
}

func New(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{c: minioClient}, err
}

func (m *MinioClient) Put(
	ctx context.Context,
	bucketName,
	objectName string,
	data []byte,
	size int64,
	contentType string,
) (info minio.UploadInfo, err error) {
	return m.c.PutObject(
		ctx,
		bucketName,
		objectName,
		bytes.NewBuffer(data),
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
}

func (m *MinioClient) Compose(
	ctx context.Context,
	bucketName,
	objectName string,
	totalChunks int,
) (info minio.UploadInfo, err error) {
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: objectName,
	}
	var srcs []minio.CopySrcOptions
	for i := 1; i <= totalChunks; i++ {
		fInt := strconv.FormatInt(int64(i), 10)
		src := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: objectName + "_" + fInt,
		}
		srcs = append(srcs, src)
	}
	return m.c.ComposeObject(
		ctx,
		dst,
		srcs...,
	)
}
