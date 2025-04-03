package test

import (
	"github.com/cloudnationhq/az-cn-go-markparsr"
	"testing"
)

func TestReadmeValidation(t *testing.T) {

	// Use functional options pattern
	validator, err := markparsr.NewReadmeValidator(
		markparsr.WithRelativeReadmePath("../module/README.md"),
		markparsr.WithAdditionalSections("Goals", "Testing", "Notes"),
		markparsr.WithAdditionalFiles("GOALS.md", "TESTING.md"),
		markparsr.WithProviderPrefixes("azurerm_", "random_", "tls_"),
	)

	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	errors := validator.Validate()
	if len(errors) > 0 {
		for _, err := range errors {
			t.Errorf("Validation error: %v", err)
		}
	}
}
