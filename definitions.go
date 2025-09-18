package markparsr

// TerraformDefinitionValidator ensures Terraform resources and data sources are documented in markdown.
type TerraformDefinitionValidator struct {
	markdown  *MarkdownContent
	terraform *TerraformContent
}

// NewTerraformDefinitionValidator creates a validator for Terraform resources and data sources.
func NewTerraformDefinitionValidator(markdown *MarkdownContent, terraform *TerraformContent) *TerraformDefinitionValidator {
	return &TerraformDefinitionValidator{
		markdown:  markdown,
		terraform: terraform,
	}
}

// Validate compares Terraform resources and data sources with markdown documentation.
func (tdv *TerraformDefinitionValidator) Validate() []error {
	tfResources, tfDataSources, err := tdv.terraform.ExtractResourcesAndDataSources()
	if err != nil {
		return []error{err}
	}

	readmeResources, readmeDataSources, mdErr := tdv.markdown.ExtractResourcesAndDataSources()

	collector := &ErrorCollector{}
	collector.Add(mdErr)
	if tdv.markdown.HasSection("Resources") || len(readmeResources) > 0 || len(readmeDataSources) > 0 {
		collector.AddMany(compareTerraformAndMarkdown(tfResources, readmeResources, "Resources"))
		collector.AddMany(compareTerraformAndMarkdown(tfDataSources, readmeDataSources, "Data Sources"))
	}

	return collector.Errors()
}
