package funcs

type Colony struct {
	rooms     map[string]*Room
	antCount  int
	roomPaths map[*Room]int
	workers   []*Ant
	startRoom *Room
	endRoom   *Room
}
type Ant struct {
}

type Room struct {
	tunnels       *TunnelList
	isStart       bool
	isDestination bool
	isUnoccupied  bool
	xCoord        int
	yCoord        int
	roomName      string
	accessMap     map[string]bool
	isEnd         bool
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
