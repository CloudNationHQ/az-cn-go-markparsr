package markparsr

import "github.com/hashicorp/hcl/v2"

// Core interfaces for the markparsr package

// FileReader interface for file operations
type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

// HCLParser interface for HCL parsing operations
type HCLParser interface {
	ParseHCL(content []byte, filename string) (*hcl.File, hcl.Diagnostics)
}

// ResourceExtractor interface for extracting Terraform resources
type ResourceExtractor interface {
	ExtractResourcesAndDataSources() ([]string, []string, error)
	ExtractItems(filePath, blockType string) ([]string, error)
}

// DocumentParser interface for basic document operations
type DocumentParser interface {
	GetContent() string
	GetAllSections() []string
	HasSection(sectionName string) bool
}

// SectionExtractor interface for extracting items from sections
type SectionExtractor interface {
	ExtractSectionItems(sectionNames ...string) []string
}

// ResourceDocumentExtractor interface for extracting resources from documentation
type ResourceDocumentExtractor interface {
	ExtractResourcesAndDataSources() ([]string, []string, error)
}

// ComparisonValidator interface for validating differences between Terraform and markdown
type ComparisonValidator interface {
	ValidateItems(tfItems, mdItems []string, itemType string) []error
}

// StringUtils interface for string utility operations
type StringUtils interface {
	LevenshteinDistance(s1, s2 string) int
	IsSimilarSection(found, expected string) bool
}

// Validator interface for validation operations
type Validator interface {
	Validate() []error
}