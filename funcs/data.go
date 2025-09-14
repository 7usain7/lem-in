package funcs

type Colony struct {
	rooms     map[string]*Room
	antCount  int
	roomPaths map[*Room]int
	ants      []*Ant
	startRoom *Room
	endRoom   *Room
}
type Ant struct {
	antNumber        int
	currentRoom      *Room
	visitedRoom      map[*Room]bool
	inMotion         bool
	hasCompletedMove bool
}

type Room struct {
	tunnels      *TunnelList
	isStart      bool
	isUnoccupied bool
	xCoord       int
	yCoord       int
	roomName     string
	accessMap    map[string]bool
	isEnd        bool
}

// tunnel between rooms
type Tunnel struct {
	data           *Room
	nextConnection *Tunnel
}

// linked list of tunnels
type TunnelList struct {
	firstNode *Tunnel
}
