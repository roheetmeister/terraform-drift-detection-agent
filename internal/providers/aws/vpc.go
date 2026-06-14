package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchVPC(ctx context.Context, cfg awscfg.Config, vpcID, region string) (*normalizer.Resource, error) {
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeVpcs(%s): %w", vpcID, err)
	}
	if len(out.Vpcs) == 0 {
		return nil, fmt.Errorf("VPC %s not found", vpcID)
	}

	vpc := out.Vpcs[0]
	return &normalizer.Resource{
		ID:       vpcID,
		Type:     "aws_vpc",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"cidr_block":       aws.ToString(vpc.CidrBlock),
			"instance_tenancy": string(vpc.InstanceTenancy),
			"is_default":       aws.ToBool(vpc.IsDefault),
			"owner_id":         aws.ToString(vpc.OwnerId),
		},
		Tags: extractTags(vpc.Tags),
	}, nil
}

func fetchSubnet(ctx context.Context, cfg awscfg.Config, subnetID, region string) (*normalizer.Resource, error) {
	client := ec2.NewFromConfig(cfg)

	out, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		SubnetIds: []string{subnetID},
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeSubnets(%s): %w", subnetID, err)
	}
	if len(out.Subnets) == 0 {
		return nil, fmt.Errorf("subnet %s not found", subnetID)
	}

	sn := out.Subnets[0]
	return &normalizer.Resource{
		ID:       subnetID,
		Type:     "aws_subnet",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"cidr_block":                  aws.ToString(sn.CidrBlock),
			"vpc_id":                      aws.ToString(sn.VpcId),
			"availability_zone":           aws.ToString(sn.AvailabilityZone),
			"map_public_ip_on_launch":     aws.ToBool(sn.MapPublicIpOnLaunch),
			"assign_ipv6_address_on_creation": aws.ToBool(sn.AssignIpv6AddressOnCreation),
		},
		Tags: extractTags(sn.Tags),
	}, nil
}
