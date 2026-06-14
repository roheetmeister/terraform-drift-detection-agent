package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchIAMRole(ctx context.Context, cfg awscfg.Config, roleName, region string) (*normalizer.Resource, error) {
	client := iam.NewFromConfig(cfg)

	out, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return nil, fmt.Errorf("GetRole(%s): %w", roleName, err)
	}

	role := out.Role
	res := &normalizer.Resource{
		ID:       aws.ToString(role.RoleId),
		Type:     "aws_iam_role",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"name":                  aws.ToString(role.RoleName),
			"arn":                   aws.ToString(role.Arn),
			"path":                  aws.ToString(role.Path),
			"description":           aws.ToString(role.Description),
			"max_session_duration":  role.MaxSessionDuration,
			"assume_role_policy":    aws.ToString(role.AssumeRolePolicyDocument),
		},
		Tags: make(map[string]string),
	}

	// IAM tags are in a separate field
	for _, t := range role.Tags {
		res.Tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}

	return res, nil
}
