package peer

type PeerChannels struct {
	UsernameCh   chan string
	RawMessageCh chan string
}

func InitChannels() *PeerChannels {
	return &PeerChannels{
		UsernameCh:   make(chan string),
		RawMessageCh: make(chan string),
	}
}
