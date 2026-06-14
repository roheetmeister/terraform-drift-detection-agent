package normalizer

// Resource is the provider-agnostic representation of a cloud resource.
type Resource struct {
	ID         string
	Type       string            // e.g. "aws_instance"
	Provider   string            // e.g. "aws"
	Region     string
	Attributes map[string]interface{}
	Tags       map[string]string
}
