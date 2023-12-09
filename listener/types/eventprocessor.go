package types

// EventProcessor queries events from source chain and transforms them to Routerchain sdk.Msg
type EventProcessor interface {
	ProcessInboundEvents(lastQueriedBlock uint64, lastProcessedEventNonce uint64) error
}
