package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3 struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

// NewS3 creates an S3-compatible storage client (Yandex Cloud, AWS, MinIO, etc.).
// endpoint example: "https://storage.yandexcloud.net"
func NewS3(bucket, region, endpoint, accessKey, secretKey string) *S3 {
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	return &S3{
		client:    client,
		bucket:    bucket,
		publicURL: fmt.Sprintf("%s/%s", endpoint, bucket),
	}
}

// DeleteByURL deletes the object whose public URL matches this bucket. URLs
// pointing to other hosts are silently ignored.
func (s *S3) DeleteByURL(ctx context.Context, url string) error {
	prefix := s.publicURL + "/"
	if !strings.HasPrefix(url, prefix) {
		return nil // external URL — not ours, skip
	}
	key := strings.TrimPrefix(url, prefix)
	if key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// Upload puts data into the bucket under the given key and returns the public URL.
func (s *S3) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("s3 upload: %w", err)
	}
	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}
