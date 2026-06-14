package state

// TFState represents the top-level Terraform state file (v4 format).
type TFState struct {
	Version          int          `json:"version"`
	TerraformVersion string       `json:"terraform_version"`
	Serial           int64        `json:"serial"`
	Lineage          string       `json:"lineage"`
	Resources        []TFResource `json:"resources"`
}

// TFResource is a resource block in Terraform state.
type TFResource struct {
	Mode      string       `json:"mode"`      // "managed" or "data"
	Type      string       `json:"type"`      // e.g. "aws_instance"
	Name      string       `json:"name"`      // logical name in config
	Provider  string       `json:"provider"`  // provider source string
	Instances []TFInstance `json:"instances"` // usually one instance per resource
}

// TFInstance holds the actual attribute values for one resource instance.
type TFInstance struct {
	SchemaVersion int                    `json:"schema_version"`
	Attributes    map[string]interface{} `json:"attributes"`
}
