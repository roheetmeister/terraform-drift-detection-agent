package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchEC2Instance(ctx context.Context, cfg awscfg.Config, instanceID, region string) (*normalizer.Resource, error) {
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeInstances(%s): %w", instanceID, err)
	}

	for _, r := range out.Reservations {
		for _, inst := range r.Instances {
			res := &normalizer.Resource{
				ID:       instanceID,
				Type:     "aws_instance",
				Provider: "aws",
				Region:   region,
				Attributes: map[string]interface{}{
					"instance_type":    string(inst.InstanceType),
					"ami":              aws.ToString(inst.ImageId),
					"availability_zone": aws.ToString(inst.Placement.AvailabilityZone),
					"subnet_id":        aws.ToString(inst.SubnetId),
					"vpc_id":           aws.ToString(inst.VpcId),
					"private_ip":       aws.ToString(inst.PrivateIpAddress),
					"public_ip":        aws.ToString(inst.PublicIpAddress),
					"state":            string(inst.State.Name),
					"key_name":         aws.ToString(inst.KeyName),
					"ebs_optimized":    inst.EbsOptimized,
					"monitoring":       string(inst.Monitoring.State),
				},
				Tags: extractTags(inst.Tags),
			}
			return res, nil
		}
	}

	return nil, fmt.Errorf("instance %s not found", instanceID)
}
