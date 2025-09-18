package markparsr

import (
	"path/filepath"
)

// ItemValidator ensures Terraform items such as variables and outputs are documented.
type ItemValidator struct {
	markdown  *MarkdownContent
	terraform *TerraformContent
	itemType  string
	blockType string
	sections  []string
	fileName  string
}

// NewItemValidator links Terraform blocks to their expected markdown sections.
func NewItemValidator(markdown *MarkdownContent, terraform *TerraformContent, itemType, blockType string, sections []string, fileName string) *ItemValidator {
	return &ItemValidator{
		markdown:  markdown,
		terraform: terraform,
		itemType:  itemType,
		blockType: blockType,
		sections:  sections,
		fileName:  fileName,
	}
}

// Validate reports mismatches between Terraform items and markdown sections.
func (iv *ItemValidator) Validate() []error {
	filePath := filepath.Join(iv.terraform.workspace, iv.fileName)
	tfItems, err := iv.terraform.ExtractItems(filePath, iv.blockType)
	if err != nil {
		return []error{err}
	}

	sectionPresent := false
	var mdItems []string
	for _, section := range iv.sections {
		if iv.markdown.HasSection(section) {
			sectionPresent = true
		}
		mdItems = append(mdItems, iv.markdown.ExtractSectionItems(section)...)
	}

	if !sectionPresent && len(mdItems) == 0 {
		return nil
	}

	return compareTerraformAndMarkdown(tfItems, mdItems, iv.itemType)
}
