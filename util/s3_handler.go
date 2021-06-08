package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/glog"
)

type S3ObjectStorageHandler struct {
	S3Client   *s3.Client
	S3Endpoint string
}

func InitS3ObjectStorageHandler() (*S3ObjectStorageHandler, error) {
	endpoint := "https://s3.computational.bio.uni-giessen.de"

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("RegionOne"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: endpoint,
				}, nil
			})),
	)

	if err != nil {
		glog.Errorln(err.Error())
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	handler := S3ObjectStorageHandler{
		S3Client:   client,
		S3Endpoint: endpoint,
	}

	return &handler, nil
}

func (handler *S3ObjectStorageHandler) DeleteObject(bucket string, key string) error {
	_, err := handler.S3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		glog.Errorln(err.Error())
		return err
	}

	return nil
}
