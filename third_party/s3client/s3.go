package s3client

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"path/filepath"
)

const (
	keyBlockSize    = 3
	keyMaxSliceSize = 3
)

type S3Storage struct {
	bucket string
	srv    *s3.S3
}

type S3Config interface {
	GetEndpoint() string
	GetBucket() string
	GetAccessKey() string
	GetSecretKey() string
	GetRegion() string
	GetSchemes() string
}

func NewS3Client(storage S3Config) *S3Storage {
	sess := session.Must(session.NewSession())
	awsConfig := &aws.Config{}
	if storage.GetEndpoint() != "" {
		awsConfig.WithEndpoint(storage.GetEndpoint())
	}
	if storage.GetRegion() != "" {
		awsConfig.WithRegion(storage.GetRegion())
	}

	if storage.GetSecretKey() != "" && storage.GetAccessKey() != "" {
		awsConfig.WithCredentials(credentials.NewStaticCredentials(storage.GetAccessKey(), storage.GetSecretKey(), ""))
	}

	if storage.GetSchemes() == "http" {
		awsConfig.WithDisableSSL(true)
		awsConfig.WithS3ForcePathStyle(true)
	}
	srv := s3.New(sess, awsConfig)
	bucket := storage.GetBucket()
	_, err := srv.HeadBucket(&s3.HeadBucketInput{Bucket: &bucket})
	if err != nil {
		panic(err)
	}

	//exist, err := client.IsBucketExist(storage.GetBucket())
	//if err != nil {
	//	panic(err)
	//}
	//if !exist {
	//	err = client.CreateBucket(storage.GetBucket())
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	return &S3Storage{
		srv:    srv,
		bucket: bucket,
	}
}

func (storage *S3Storage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := storage.srv.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &storage.bucket,
		Key:    storage.transformKey(key),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (storage *S3Storage) Put(ctx context.Context, key string, hash string, data []byte) error {
	contentLength := int64(len(data))
	_, err := storage.srv.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Key:    storage.transformKey(key),
		Body:   bytes.NewReader(data),
		Bucket: &storage.bucket,

		ContentLength: &contentLength,
		Metadata: map[string]*string{
			"ipfs": &hash,
		},
	})
	return err
}

func (storage *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := storage.srv.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: &storage.bucket,
		Key:    storage.transformKey(key),
	})
	return err
}

func (storage *S3Storage) GetSize(ctx context.Context, key string) (int, error) {
	output, err := storage.srv.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: &storage.bucket,
		Key:    storage.transformKey(key),
	})
	if err != nil {
		return 0, err
	}
	if output.ContentLength == nil {
		return 0, errors.New("parse content length failed")
	}
	return int(*output.ContentLength), nil
}

func (storage *S3Storage) transformKey(key string) *string {
	size := keyMaxSliceSize
	if c := len(key) / keyBlockSize; c <= keyMaxSliceSize {
		size = c
	}
	suffix := key[len(key)-size*keyBlockSize:]
	pathSlice := make([]string, size+1)
	pathSlice[size] = key
	for i := 0; i < size; i++ {
		from, to := i*keyBlockSize, (i*keyBlockSize)+keyBlockSize
		pathSlice[i] = suffix[from:to]
	}
	return aws.String(filepath.Join(pathSlice...))
}
