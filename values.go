package gitconfig

type Values []Value

func (v Values) Len() int { return len([]Value(v)) }

func (v Values) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v Values) Less(i, j int) bool { return v[i].Name() < v[j].Name() }
