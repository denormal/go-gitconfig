package gitconfig

type Properties []Property

func (p Properties) Len() int { return len([]Property(p)) }

func (p Properties) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p Properties) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
