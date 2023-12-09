package transmitter

type IChainTransmitter interface {
	AddHandleRequestToMsgChannel(messages []byte)
	DestinationChainId() string
}
