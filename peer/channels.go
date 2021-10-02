package peer

type PeerChannels struct {
	UsernameCh   chan string
	InvalidCh    chan bool
	RawMessageCh chan string
}

func InitChannels() *PeerChannels {
	return &PeerChannels{
		UsernameCh:   make(chan string),
		InvalidCh:    make(chan bool),
		RawMessageCh: make(chan string),
	}
}
