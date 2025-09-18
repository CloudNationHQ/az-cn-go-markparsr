package markparsr

import "fmt"

// ValidationSeverity represents the severity level of a validation issue
type ValidationSeverity int

const (
	// SeverityError represents a critical validation failure
	SeverityError ValidationSeverity = iota
	// SeverityWarning represents a non-critical validation issue
	SeverityWarning
)

// ValidationIssue represents a single validation problem
type ValidationIssue struct {
	Message  string
	Severity ValidationSeverity
	Source   string // File or section where the issue was found
}

// ValidationResult holds the results of a validation operation
type ValidationResult struct {
	Errors   []ValidationIssue
	Warnings []ValidationIssue
}

// AddError adds an error-level validation issue
func (vr *ValidationResult) AddError(message, source string) {
	vr.Errors = append(vr.Errors, ValidationIssue{
		Message:  message,
		Severity: SeverityError,
		Source:   source,
	})
}

// AddWarning adds a warning-level validation issue
func (vr *ValidationResult) AddWarning(message, source string) {
	vr.Warnings = append(vr.Warnings, ValidationIssue{
		Message:  message,
		Severity: SeverityWarning,
		Source:   source,
	})
}

// HasErrors returns true if there are any error-level issues
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are any warning-level issues
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// AllIssues returns all issues (errors and warnings) combined
func (vr *ValidationResult) AllIssues() []ValidationIssue {
	var all []ValidationIssue
	all = append(all, vr.Errors...)
	all = append(all, vr.Warnings...)
	return all
}

// ToErrors converts validation result to a slice of errors for backward compatibility
func (vr *ValidationResult) ToErrors() []error {
	var errors []error
	for _, issue := range vr.Errors {
		errors = append(errors, fmt.Errorf("%s", issue.Message))
	}
	for _, issue := range vr.Warnings {
		errors = append(errors, fmt.Errorf("%s", issue.Message))
	}
	return errors
}

// ErrorCollector accumulates errors while skipping nil entries.
type ErrorCollector struct {
	errs []error
}

// Add records a single error when it is non-nil.
func (c *ErrorCollector) Add(err error) {
	if err != nil {
		c.errs = append(c.errs, err)
	}
}

// AddMany records a slice of errors, ignoring nil entries.
func (c *ErrorCollector) AddMany(errs []error) {
	for _, err := range errs {
		c.Add(err)
	}
}

// Errors returns the collected errors.
func (c *ErrorCollector) Errors() []error {
	return c.errs
}

// HasErrors reports whether any errors were collected.
func (c *ErrorCollector) HasErrors() bool {
	return len(c.errs) > 0
}
