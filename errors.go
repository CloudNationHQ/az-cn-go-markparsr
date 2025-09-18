package markparsr

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
