// Package markparsr provides utilities for validating Terraform module documentation.
package markparsr

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownFormat string

const (
	FormatDocument MarkdownFormat = "document"
)

var (
	inputAnchorRe  = regexp.MustCompile(`(?i)<a\s+name="input_([^"\s]+)"`)
	outputAnchorRe = regexp.MustCompile(`(?i)<a\s+name="output_([^"\s]+)"`)
)

// MarkdownContent parses and analyzes Terraform module documentation
type MarkdownContent struct {
	data             string
	rootNode         ast.Node
	sections         map[string]bool
	format           MarkdownFormat
	stringPool       *sync.Pool
	providerPrefixes []string
	h2Headings       []*ast.Heading
	sectionNames     []string
	sectionMatches   map[string][]*ast.Heading
	anchorTypes      map[string]map[string]bool
}

// NewMarkdownContent creates a new analyzer for markdown content
func NewMarkdownContent(data string, format MarkdownFormat, providerPrefixes []string) *MarkdownContent {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	rootNode := markdown.Parse([]byte(data), p)

	mc := &MarkdownContent{
		data:     data,
		rootNode: rootNode,
		sections: make(map[string]bool),
		stringPool: &sync.Pool{
			New: func() any {
				return &strings.Builder{}
			},
		},
		providerPrefixes: providerPrefixes,
		sectionMatches:   make(map[string][]*ast.Heading),
	}

	mc.indexHeadings()
	mc.indexAnchors()

	mc.format = FormatDocument
	if format != "" && format != FormatDocument {
		fmt.Printf("Markdown format '%s' is not supported; using document format\n", format)
	}

	return mc
}

func (mc *MarkdownContent) indexHeadings() {
	var names []string
	ast.WalkFunc(mc.rootNode, func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}
		if heading, ok := node.(*ast.Heading); ok && heading.Level == 2 {
			mc.h2Headings = append(mc.h2Headings, heading)
			text := strings.TrimSpace(mc.extractText(heading))
			if text != "" {
				names = append(names, text)
			}
			return ast.SkipChildren
		}
		return ast.GoToNext
	})
	mc.sectionNames = names
}

func (mc *MarkdownContent) indexAnchors() {
	mc.anchorTypes = make(map[string]map[string]bool)
	type anchorDef struct {
		re  *regexp.Regexp
		typ string
	}
	defs := []anchorDef{
		{re: inputAnchorRe, typ: "input"},
		{re: outputAnchorRe, typ: "output"},
	}

	for _, def := range defs {
		matches := def.re.FindAllStringSubmatch(mc.data, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			name := strings.ToLower(strings.TrimSpace(match[1]))
			if name == "" {
				continue
			}
			if mc.anchorTypes[name] == nil {
				mc.anchorTypes[name] = make(map[string]bool)
			}
			mc.anchorTypes[name][def.typ] = true
		}
	}
}

// GetContent returns the full markdown content
func (mc *MarkdownContent) GetContent() string {
	return mc.data
}

func (mc *MarkdownContent) HasSection(sectionName string) bool {
	if found, exists := mc.sections[sectionName]; exists {
		return found
	}

	found := len(mc.matchSectionHeadings(sectionName)) > 0
	mc.sections[sectionName] = found
	return found
}

// GetAllSections returns a list of all H2 section names in the markdown
func (mc *MarkdownContent) GetAllSections() []string {
	return append([]string(nil), mc.sectionNames...)
}

// ExtractSectionItems extracts item names from a section using document-style headings
func (mc *MarkdownContent) ExtractSectionItems(sectionNames ...string) []string {
	return mc.extractDocumentSectionItems(sectionNames...)
}

func (mc *MarkdownContent) matchSectionHeadings(sectionName string) []*ast.Heading {
	key := strings.ToLower(strings.TrimSpace(sectionName))
	if key == "" {
		return nil
	}

	if cached, ok := mc.sectionMatches[key]; ok {
		return cached
	}

	var matches []*ast.Heading
	for _, heading := range mc.h2Headings {
		text := strings.TrimSpace(mc.extractText(heading))
		if matchesSectionName(text, sectionName) {
			matches = append(matches, heading)
		}
	}

	mc.sectionMatches[key] = matches
	return matches
}

func (mc *MarkdownContent) collectSectionHeadings(sectionNames []string) []*ast.Heading {
	if len(sectionNames) == 0 {
		return nil
	}

	seen := make(map[*ast.Heading]struct{})
	var headings []*ast.Heading
	for _, name := range sectionNames {
		for _, heading := range mc.matchSectionHeadings(name) {
			if _, ok := seen[heading]; ok {
				continue
			}
			seen[heading] = struct{}{}
			headings = append(headings, heading)
		}
	}

	return headings
}

// extractDocumentSectionItems extracts item names from level 3 headings within specified sections
func (mc *MarkdownContent) extractDocumentSectionItems(sectionNames ...string) []string {
	headings := mc.collectSectionHeadings(sectionNames)
	if len(headings) == 0 {
		return mc.fallbackSectionItems(sectionNames)
	}

	var items []string
	for _, heading := range headings {
		items = append(items, mc.itemsUnderHeading(heading)...)
	}

	return mc.filterItemsByAnchorType(sectionNames, items)
}

func (mc *MarkdownContent) fallbackSectionItems(sectionNames []string) []string {
	needInputs := false
	needOutputs := false
	for _, name := range sectionNames {
		lower := strings.ToLower(name)
		if strings.Contains(lower, "input") {
			needInputs = true
		}
		if strings.Contains(lower, "output") {
			needOutputs = true
		}
	}

	var items []string
	seen := make(map[string]struct{})

	if needInputs {
		for _, item := range mc.extractAnchoredItems(inputAnchorRe) {
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			items = append(items, item)
		}
	}

	if needOutputs {
		for _, item := range mc.extractAnchoredItems(outputAnchorRe) {
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			items = append(items, item)
		}
	}

	return mc.filterItemsByAnchorType(sectionNames, items)
}

func (mc *MarkdownContent) filterItemsByAnchorType(sectionNames, items []string) []string {
	if len(items) == 0 {
		return items
	}

	expected := mc.expectedAnchorType(sectionNames)
	if expected == "" {
		return items
	}

	var filtered []string
	for _, item := range items {
		anchorSet := mc.anchorTypes[strings.ToLower(item)]
		if len(anchorSet) > 0 && !anchorSet[expected] {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered
}

func (mc *MarkdownContent) expectedAnchorType(sectionNames []string) string {
	for _, name := range sectionNames {
		lower := strings.ToLower(name)
		if strings.Contains(lower, "input") {
			return "input"
		}
		if strings.Contains(lower, "output") {
			return "output"
		}
	}
	return ""
}

func (mc *MarkdownContent) extractAnchoredItems(re *regexp.Regexp) []string {
	matches := re.FindAllStringSubmatch(mc.data, -1)
	var items []string
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := strings.TrimSpace(match[1])
		if name == "" {
			continue
		}
		items = append(items, name)
	}
	return items
}

func (mc *MarkdownContent) itemsUnderHeading(heading *ast.Heading) []string {
	var items []string
	for node := getNextSibling(heading); node != nil; node = getNextSibling(node) {
		if h, ok := node.(*ast.Heading); ok {
			if h.Level <= heading.Level {
				break
			}
			if h.Level == 3 {
				if name, ok := mc.itemNameFromHeading(h); ok {
					items = append(items, name)
				}
			}
			continue
		}
	}
	return items
}

func (mc *MarkdownContent) itemNameFromHeading(heading *ast.Heading) (string, bool) {
	headingText := strings.TrimSpace(mc.extractText(heading))
	if headingText == "" {
		return "", false
	}

	name := strings.Trim(headingText, " []")
	name = strings.TrimPrefix(name, "<a name=\"input_")
	name = strings.TrimPrefix(name, "<a name=\"output_")
	name = strings.TrimSuffix(name, "</a>")
	name = strings.TrimSuffix(name, "\"></a>")
	name = strings.TrimSpace(name)
	if name == "" {
		return "", false
	}
	return name, true
}

func (mc *MarkdownContent) resourcesUnderHeading(heading *ast.Heading) ([]string, []string) {
	var resources []string
	var dataSources []string

	for node := getNextSibling(heading); node != nil; node = getNextSibling(node) {
		if h, ok := node.(*ast.Heading); ok && h.Level <= heading.Level {
			break
		}

		ast.WalkFunc(node, func(n ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.GoToNext
			}
			if link, ok := n.(*ast.Link); ok {
				mc.appendResourceFromLink(link, &resources, &dataSources)
				return ast.SkipChildren
			}
			return ast.GoToNext
		})
	}

	return resources, dataSources
}

func (mc *MarkdownContent) resourcesWithoutHeading() ([]string, []string) {
	var resources []string
	var dataSources []string

	ast.WalkFunc(mc.rootNode, func(n ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}
		if link, ok := n.(*ast.Link); ok {
			mc.appendResourceFromLink(link, &resources, &dataSources)
			return ast.SkipChildren
		}
		return ast.GoToNext
	})

	return resources, dataSources
}

func (mc *MarkdownContent) appendResourceFromLink(link *ast.Link, resources, dataSources *[]string) {
	linkText := mc.extractText(link)
	destination := string(link.Destination)
	if !mc.hasProviderPrefix(linkText) {
		return
	}

	resourceName := strings.Split(linkText, "]")[0]
	resourceName = strings.TrimPrefix(resourceName, "[")
	baseName := strings.Split(resourceName, ".")[0]

	if strings.Contains(destination, "/data-sources/") {
		addUnique(dataSources, resourceName)
		addUnique(dataSources, baseName)
	} else {
		addUnique(resources, resourceName)
		addUnique(resources, baseName)
	}
}

// ExtractResourcesAndDataSources finds Terraform resources and data sources in the markdown
func (mc *MarkdownContent) ExtractResourcesAndDataSources() ([]string, []string, error) {
	return mc.extractDocumentResourcesAndDataSources()
}

// extractDocumentResourcesAndDataSources finds resources in document style markdown
func (mc *MarkdownContent) extractDocumentResourcesAndDataSources() ([]string, []string, error) {
	headings := mc.collectSectionHeadings([]string{"Resources"})
	var resources []string
	var dataSources []string

	if len(headings) == 0 {
		resources, dataSources = mc.resourcesWithoutHeading()
		if len(resources) == 0 && len(dataSources) == 0 {
			return nil, nil, fmt.Errorf("resources section not found or empty")
		}
		return resources, dataSources, fmt.Errorf("resources section not found or empty")
	}

	for _, heading := range headings {
		r, d := mc.resourcesUnderHeading(heading)
		resources = append(resources, r...)
		dataSources = append(dataSources, d...)
	}

	if len(resources) == 0 && len(dataSources) == 0 {
		return nil, nil, fmt.Errorf("resources section not found or empty")
	}

	return resources, dataSources, nil
}

// extractText gets the text content from a node, using a string pool for efficiency
func (mc *MarkdownContent) extractText(node ast.Node) string {
	sb := mc.stringPool.Get().(*strings.Builder)
	sb.Reset()
	defer mc.stringPool.Put(sb)

	ast.WalkFunc(node, func(n ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch tn := n.(type) {
			case *ast.Text:
				sb.Write(tn.Literal)
			case *ast.Code:
				sb.Write(tn.Literal)
			}
		}
		return ast.GoToNext
	})

	return sb.String()
}

// hasProviderPrefix checks if a string has a recognized provider prefix
func (mc *MarkdownContent) hasProviderPrefix(s string) bool {
	s = strings.ToLower(s)

	// If no provider prefixes are configured, use default behavior
	if len(mc.providerPrefixes) == 0 {
		return false
	}

	for _, prefix := range mc.providerPrefixes {
		if strings.HasPrefix(s, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// getNextSibling returns the next sibling of a node.
func getNextSibling(node ast.Node) ast.Node {
	parent := node.GetParent()
	if parent == nil {
		return nil
	}
	children := parent.GetChildren()
	for i, n := range children {
		if n == node && i+1 < len(children) {
			return children[i+1]
		}
	}
	return nil
}

// addUnique adds a string to a slice if it's not already present
func addUnique(slice *[]string, item string) {
	if !slices.Contains(*slice, item) {
		*slice = append(*slice, item)
	}
}

func matchesSectionName(actual, expected string) bool {
	actual = strings.TrimSpace(actual)
	expected = strings.TrimSpace(expected)

	if actual == "" || expected == "" {
		return false
	}

	if strings.EqualFold(actual, expected) {
		return true
	}

	if strings.EqualFold(actual, expected+"s") || strings.EqualFold(actual+"s", expected) {
		return true
	}

	if expected == "Inputs" && (strings.EqualFold(actual, "Required Inputs") || strings.EqualFold(actual, "Optional Inputs")) {
		return true
	}

	return isSimilarSection(actual, expected)
}
