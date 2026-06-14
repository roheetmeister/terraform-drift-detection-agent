package aws

import (
	"context"
	"fmt"
	"strings"

	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchS3Bucket(ctx context.Context, cfg awscfg.Config, bucketName, region string) (*normalizer.Resource, error) {
	client := s3.NewFromConfig(cfg)

	locOut, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: awscfg.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("GetBucketLocation(%s): %w", bucketName, err)
	}

	bucketRegion := string(locOut.LocationConstraint)
	if bucketRegion == "" {
		bucketRegion = "us-east-1" // empty LocationConstraint means us-east-1
	}

	res := &normalizer.Resource{
		ID:       bucketName,
		Type:     "aws_s3_bucket",
		Provider: "aws",
		Region:   bucketRegion,
		Attributes: map[string]interface{}{
			"bucket": bucketName,
			"region": bucketRegion,
		},
		Tags: make(map[string]string),
	}

	// Fetch bucket tags
	tagOut, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: awscfg.String(bucketName),
	})
	if err == nil {
		for _, t := range tagOut.TagSet {
			res.Tags[awscfg.ToString(t.Key)] = awscfg.ToString(t.Value)
		}
	} else if !strings.Contains(err.Error(), "NoSuchTagSet") {
		// NoSuchTagSet is expected for untagged buckets — ignore it
		_ = err
	}

	// Fetch versioning status
	verOut, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: awscfg.String(bucketName),
	})
	if err == nil {
		res.Attributes["versioning_enabled"] = string(verOut.Status) == "Enabled"
	}

	return res, nil
}
