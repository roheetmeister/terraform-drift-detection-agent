package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchSecurityGroup(ctx context.Context, cfg awscfg.Config, sgID, region string) (*normalizer.Resource, error) {
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		GroupIds: []string{sgID},
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeSecurityGroups(%s): %w", sgID, err)
	}

	if len(out.SecurityGroups) == 0 {
		return nil, fmt.Errorf("security group %s not found", sgID)
	}

	sg := out.SecurityGroups[0]
	return &normalizer.Resource{
		ID:       sgID,
		Type:     "aws_security_group",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"name":        aws.ToString(sg.GroupName),
			"description": aws.ToString(sg.Description),
			"vpc_id":      aws.ToString(sg.VpcId),
			"owner_id":    aws.ToString(sg.OwnerId),
		},
		Tags: extractTags(sg.Tags),
	}, nil
}
