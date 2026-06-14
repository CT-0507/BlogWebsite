package infrastructure

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service struct {
	client           *s3.Client
	bucket           string
	cloudFrontDomain string
}

const (
	TEMPORARY_FLAG = "temporary"
)

func NewS3Service(
	client *s3.Client,
	bucket string,
	cloudFrontDomain string,
) *S3Service {
	return &S3Service{
		client:           client,
		bucket:           bucket,
		cloudFrontDomain: strings.TrimRight(cloudFrontDomain, "/"),
	}
}

func (s *S3Service) buildURL(key string) string {
	key = strings.TrimLeft(key, "/")
	return fmt.Sprintf("%s/%s", s.cloudFrontDomain, key)
}
func (s *S3Service) Save(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
	isTemporary bool,
) (*storage.UploadResult, error) {

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		Tagging:     aws.String(fmt.Sprintf("%s=%t", TEMPORARY_FLAG, isTemporary)),
	})

	if err != nil {
		return nil, err
	}

	return &storage.UploadResult{
		Key: key,
		URL: s.buildURL(key),
	}, nil
}

func (s *S3Service) MarkPermanent(
	ctx context.Context,
	key string,
) error {

	return s.changeFileTemporaryFlag(ctx, key, false)
}

func (s *S3Service) MarkDelete(
	ctx context.Context,
	key string,
) error {

	return s.changeFileTemporaryFlag(ctx, key, true)
}

func (s *S3Service) changeFileTemporaryFlag(
	ctx context.Context,
	key string,
	value bool,
) error {
	_, err := s.client.PutObjectTagging(
		ctx,
		&s3.PutObjectTaggingInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
			Tagging: &types.Tagging{
				TagSet: []types.Tag{
					{
						Key:   aws.String("temporary"),
						Value: aws.String(strconv.FormatBool(value)),
					},
				},
			},
		},
	)

	return err
}
func (s *S3Service) Delete(
	ctx context.Context,
	key string,
) error {

	_, err := s.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)

	return err
}
