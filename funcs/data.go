package funcs

type Colony struct {
	roomMap     map[string]*Room
	antCount    int
	roomPaths   map[*Room]int
	workers     []*Ant
	initialRoom *Room
	finalRoom   *Room
}
type Ant struct {
}

type Room struct {
}
