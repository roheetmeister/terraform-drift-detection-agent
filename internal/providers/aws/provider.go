package aws

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

// Provider implements providers.Provider for AWS.
type Provider struct{}

// New creates a new AWS provider.
func New() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string { return "aws" }

func (p *Provider) SupportedTypes() []string {
	return []string{
		"aws_instance",
		"aws_s3_bucket",
		"aws_security_group",
		"aws_vpc",
		"aws_subnet",
		"aws_iam_role",
		"aws_lambda_function",
		"aws_db_instance",
		"aws_rds_cluster",
	}
}

// FetchResource dispatches to the correct AWS service fetcher by resource type.
func (p *Provider) FetchResource(ctx context.Context, resourceType, resourceID, region string) (*normalizer.Resource, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	switch resourceType {
	case "aws_instance":
		return fetchEC2Instance(ctx, cfg, resourceID, region)
	case "aws_s3_bucket":
		return fetchS3Bucket(ctx, cfg, resourceID, region)
	case "aws_security_group":
		return fetchSecurityGroup(ctx, cfg, resourceID, region)
	case "aws_vpc":
		return fetchVPC(ctx, cfg, resourceID, region)
	case "aws_subnet":
		return fetchSubnet(ctx, cfg, resourceID, region)
	case "aws_iam_role":
		return fetchIAMRole(ctx, cfg, resourceID, region)
	case "aws_lambda_function":
		return fetchLambdaFunction(ctx, cfg, resourceID, region)
	case "aws_db_instance":
		return fetchDBInstance(ctx, cfg, resourceID, region)
	case "aws_rds_cluster":
		return fetchRDSCluster(ctx, cfg, resourceID, region)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}
