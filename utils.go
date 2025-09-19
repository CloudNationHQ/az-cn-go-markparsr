package markparsr

// defaultStringUtils provides the package-default StringUtils implementation.
type defaultStringUtils struct{}

// NewStringUtils returns the default StringUtils implementation.
func NewStringUtils() StringUtils {
	return &defaultStringUtils{}
}

// LevenshteinDistance calculates the edit distance between two strings
func (dsu *defaultStringUtils) LevenshteinDistance(s1, s2 string) int {
	return levenshtein(s1, s2)
}

// IsSimilarSection checks if a found section name is likely a typo of an expected section
func (dsu *defaultStringUtils) IsSimilarSection(found, expected string) bool {
	return isSimilarSection(found, expected)
}

// levenshtein calculates the edit distance between two strings
func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	v0 := make([]int, len(s2)+1)
	v1 := make([]int, len(s2)+1)

	for i := range v0 {
		v0[i] = i
	}

	for i := range s1 {
		v1[0] = i + 1
		for j := range s2 {
			cost := 1
			if s1[i] == s2[j] {
				cost = 0
			}
			v1[j+1] = min(v1[j]+1, v0[j+1]+1, v0[j]+cost)
		}
		copy(v0, v1)
	}
	return v1[len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
