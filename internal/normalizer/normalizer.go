package normalizer

import (
	"strings"

	"github.com/roheetmeister/terraform-drift-detection-agent/internal/state"
)

// FromState converts a Terraform state resource instance into a common Resource.
func FromState(res *state.TFResource, inst *state.TFInstance, defaultRegion string) *Resource {
	r := &Resource{
		Type:       res.Type,
		Provider:   extractProvider(res.Provider),
		Region:     defaultRegion,
		Attributes: make(map[string]interface{}),
		Tags:       make(map[string]string),
	}

	for k, v := range inst.Attributes {
		r.Attributes[k] = v
	}

	if id, ok := inst.Attributes["id"].(string); ok {
		r.ID = id
	}
	if region, ok := inst.Attributes["region"].(string); ok && region != "" {
		r.Region = region
	}

	// Extract tags from the "tags" attribute (map[string]interface{})
	if tags, ok := inst.Attributes["tags"].(map[string]interface{}); ok {
		for k, v := range tags {
			if s, ok := v.(string); ok {
				r.Tags[k] = s
			}
		}
	}

	return r
}

// extractProvider derives a short provider name from the TF provider source string.
// e.g. `provider["registry.terraform.io/hashicorp/aws"]` → "aws"
func extractProvider(provider string) string {
	parts := strings.Split(provider, "/")
	if len(parts) == 0 {
		return "unknown"
	}
	last := parts[len(parts)-1]
	return strings.Trim(last, `"`)
}
