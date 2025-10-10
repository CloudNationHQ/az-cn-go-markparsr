package markparsr

import (
	"path/filepath"
)

type ItemValidator struct {
	markdown  *MarkdownContent
	terraform *TerraformContent
	itemType  string
	blockType string
	sections  []string
	fileName  string
}

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
