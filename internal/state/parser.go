package state

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Parse reads a Terraform state file from a local path or an S3 URI (s3://bucket/key).
func Parse(ctx context.Context, path string) (*TFState, error) {
	if strings.HasPrefix(path, "s3://") {
		return parseS3(ctx, path)
	}
	return parseLocal(path)
}

func parseLocal(path string) (*TFState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file %q: %w", path, err)
	}
	return decode(data)
}

func parseS3(ctx context.Context, uri string) (*TFState, error) {
	trimmed := strings.TrimPrefix(uri, "s3://")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid S3 URI %q: expected s3://bucket/key", uri)
	}
	bucket, key := parts[0], parts[1]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	out, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching state from S3 %q: %w", uri, err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("reading S3 object body: %w", err)
	}

	return decode(data)
}

func decode(data []byte) (*TFState, error) {
	var st TFState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("parsing state JSON: %w", err)
	}
	if st.Version != 4 {
		return nil, fmt.Errorf("unsupported state version %d (only v4 is supported)", st.Version)
	}
	return &st, nil
}
