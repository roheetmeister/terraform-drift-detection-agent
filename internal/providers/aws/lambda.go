package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

func fetchLambdaFunction(ctx context.Context, cfg awscfg.Config, functionName, region string) (*normalizer.Resource, error) {
	client := lambda.NewFromConfig(cfg)

	out, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return nil, fmt.Errorf("GetFunction(%s): %w", functionName, err)
	}

	fn := out.Configuration
	res := &normalizer.Resource{
		ID:       aws.ToString(fn.FunctionArn),
		Type:     "aws_lambda_function",
		Provider: "aws",
		Region:   region,
		Attributes: map[string]interface{}{
			"function_name": aws.ToString(fn.FunctionName),
			"runtime":       string(fn.Runtime),
			"handler":       aws.ToString(fn.Handler),
			"role":          aws.ToString(fn.Role),
			"memory_size":   fn.MemorySize,
			"timeout":       fn.Timeout,
			"description":   aws.ToString(fn.Description),
			"architectures": fn.Architectures,
		},
		Tags: make(map[string]string),
	}

	// Lambda tags come from the Tags field on the GetFunction response
	for k, v := range out.Tags {
		res.Tags[k] = v
	}

	return res, nil
}
