package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchDBInstance(ctx context.Context, cfg awscfg.Config, dbID, region string) (*normalizer.Resource, error) {
	client := rds.NewFromConfig(cfg)

	out, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbID),
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeDBInstances(%s): %w", dbID, err)
	}
	if len(out.DBInstances) == 0 {
		return nil, fmt.Errorf("DB instance %s not found", dbID)
	}

	db := out.DBInstances[0]
	res := &normalizer.Resource{
		ID:       aws.ToString(db.DBInstanceIdentifier),
		Type:     "aws_db_instance",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"engine":               aws.ToString(db.Engine),
			"engine_version":       aws.ToString(db.EngineVersion),
			"instance_class":       aws.ToString(db.DBInstanceClass),
			"allocated_storage":    db.AllocatedStorage,
			"storage_type":         aws.ToString(db.StorageType),
			"multi_az":             db.MultiAZ,
			"publicly_accessible":  db.PubliclyAccessible,
			"deletion_protection": db.DeletionProtection,
			"db_name":              aws.ToString(db.DBName),
		},
		Tags: make(map[string]string),
	}

	for _, t := range db.TagList {
		res.Tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}

	return res, nil
}

func fetchRDSCluster(ctx context.Context, cfg awscfg.Config, clusterID, region string) (*normalizer.Resource, error) {
	client := rds.NewFromConfig(cfg)

	out, err := client.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterID),
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeDBClusters(%s): %w", clusterID, err)
	}
	if len(out.DBClusters) == 0 {
		return nil, fmt.Errorf("RDS cluster %s not found", clusterID)
	}

	cl := out.DBClusters[0]
	res := &normalizer.Resource{
		ID:       aws.ToString(cl.DBClusterIdentifier),
		Type:     "aws_rds_cluster",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"engine":               aws.ToString(cl.Engine),
			"engine_version":       aws.ToString(cl.EngineVersion),
			"database_name":        aws.ToString(cl.DatabaseName),
			"master_username":      aws.ToString(cl.MasterUsername),
			"deletion_protection": cl.DeletionProtection,
			"storage_encrypted":    cl.StorageEncrypted,
		},
		Tags: make(map[string]string),
	}

	for _, t := range cl.TagList {
		res.Tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}

	return res, nil
}
