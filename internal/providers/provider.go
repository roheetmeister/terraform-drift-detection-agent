package providers

import (
	"context"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/normalizer"
)

// Provider is the interface every cloud provider adapter must implement.
// Adding GCP/Azure = implement this interface and register in the provider registry.
type Provider interface {
	// Name returns the short provider identifier, e.g. "aws".
	Name() string

	// SupportedTypes lists the Terraform resource types this provider can fetch.
	SupportedTypes() []string

	// FetchResource retrieves live metadata for a resource by its Terraform resource type and ID.
	// Returns nil, error if the resource cannot be found or the API call fails.
	FetchResource(ctx context.Context, resourceType, resourceID, region string) (*normalizer.Resource, error)
}
