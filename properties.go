package gitconfig

// properties is the sort class for configuration properties
type properties []Property

// Len returns the length of the list of properties
func (p properties) Len() int { return len([]Property(p)) }

// Swap interchanges the position of properties in positions i and j.
func (p properties) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less returns true is the property at position i should appear earlier in
// the sort order than the property at position j.
func (p properties) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
