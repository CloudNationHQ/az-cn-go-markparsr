package markparsr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)


// defaultFileReader implements FileReader using os package
type defaultFileReader struct{}

func (dfr *defaultFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// pooledHCLParser implements HCLParser using a sync.Pool
type pooledHCLParser struct {
	pool *sync.Pool
}

func newPooledHCLParser() *pooledHCLParser {
	return &pooledHCLParser{
		pool: &sync.Pool{
			New: func() any {
				return hclparse.NewParser()
			},
		},
	}
}

func (php *pooledHCLParser) ParseHCL(content []byte, filename string) (*hcl.File, hcl.Diagnostics) {
	parser := php.pool.Get().(*hclparse.Parser)
	defer php.pool.Put(parser)
	return parser.ParseHCL(content, filename)
}

// TerraformContent extracts Terraform definitions for documentation validation.
type TerraformContent struct {
	workspace  string
	fileReader FileReader
	hclParser  HCLParser
	fileCache  sync.Map
}

// NewTerraformContent creates a Terraform analyzer rooted at the provided module path.
func NewTerraformContent(modulePath string) (*TerraformContent, error) {
	if modulePath == "" {
		githubWorkspace := os.Getenv("GITHUB_WORKSPACE")
		if githubWorkspace != "" {
			modulePath = githubWorkspace
		} else {
			var err error
			modulePath, err = os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("failed to get current working directory: %w", err)
			}
		}
	}

	return &TerraformContent{
		workspace:  modulePath,
		fileReader: &defaultFileReader{},
		hclParser:  newPooledHCLParser(),
	}, nil
}

// parseFile reads and parses an HCL file, handling common error cases
func (tc *TerraformContent) parseFile(filePath string) (*hcl.File, error) {
	content, err := tc.fileReader.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Return nil file for non-existent files
		}
		return nil, fmt.Errorf("error reading file %s: %w", filepath.Base(filePath), err)
	}

	file, parseDiags := tc.hclParser.ParseHCL(content, filePath)
	if parseDiags.HasErrors() {
		return nil, fmt.Errorf("error parsing HCL in %s: %v", filepath.Base(filePath), parseDiags)
	}

	return file, nil
}

// ExtractItems gets items of a specific block type from a Terraform file.
func (tc *TerraformContent) ExtractItems(filePath, blockType string) ([]string, error) {
	file, err := tc.parseFile(filePath)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return []string{}, nil // File doesn't exist
	}

	return tc.extractItemsFromFile(file, filePath, blockType)
}

// extractItemsFromFile extracts items of a specific block type from a parsed HCL file
func (tc *TerraformContent) extractItemsFromFile(file *hcl.File, filePath, blockType string) ([]string, error) {
	var items []string
	body := file.Body
	hclContent, _, diags := body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: blockType, LabelNames: []string{"name"}},
		},
	})

	if diags.HasErrors() {
		return nil, fmt.Errorf("error getting content from %s: %v", filepath.Base(filePath), diags)
	}

	if hclContent == nil {
		return items, nil
	}

	for _, block := range hclContent.Blocks {
		if len(block.Labels) > 0 {
			itemName := strings.TrimSpace(block.Labels[0])
			items = append(items, itemName)
		}
	}

	return items, nil
}

// ExtractResourcesAndDataSources finds resources and data sources in module-level Terraform files.
func (tc *TerraformContent) ExtractResourcesAndDataSources() ([]string, []string, error) {
	var resources []string
	var dataSources []string

	files, err := os.ReadDir(tc.workspace)
	if err != nil {
		if os.IsNotExist(err) {
			return resources, dataSources, nil
		}
		return nil, nil, fmt.Errorf("error reading directory %s: %w", tc.workspace, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}

		filePath := filepath.Join(tc.workspace, file.Name())
		fileResources, fileDataSources, err := tc.extractFromFilePath(filePath)
		if err != nil {
			return nil, nil, err
		}

		resources = append(resources, fileResources...)
		dataSources = append(dataSources, fileDataSources...)
	}

	return resources, dataSources, nil
}

// extractFromFilePath gets resources and data sources from a single Terraform file.
func (tc *TerraformContent) extractFromFilePath(filePath string) ([]string, []string, error) {
	file, err := tc.parseFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	if file == nil {
		return []string{}, []string{}, nil // File doesn't exist
	}

	return tc.extractResourcesFromFile(file, filePath)
}

// extractResourcesFromFile extracts resources and data sources from a parsed HCL file
func (tc *TerraformContent) extractResourcesFromFile(file *hcl.File, filePath string) ([]string, []string, error) {

	var resources []string
	var dataSources []string
	body := file.Body
	hclContent, _, diags := body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}},
			{Type: "data", LabelNames: []string{"type", "name"}},
		},
	})

	if diags.HasErrors() {
		return nil, nil, fmt.Errorf("error getting content from %s: %v", filepath.Base(filePath), diags)
	}

	if hclContent == nil {
		return resources, dataSources, nil
	}

	for _, block := range hclContent.Blocks {
		if len(block.Labels) >= 2 {
			resourceType := strings.TrimSpace(block.Labels[0])
			resourceName := strings.TrimSpace(block.Labels[1])
			fullResourceName := resourceType + "." + resourceName

			switch block.Type {
			case "resource":
				resources = append(resources, resourceType)
				resources = append(resources, fullResourceName)
			case "data":
				dataSources = append(dataSources, resourceType)
				dataSources = append(dataSources, fullResourceName)
			}
		}
	}

	return resources, dataSources, nil
}
