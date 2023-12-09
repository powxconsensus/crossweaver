package types

// Transformer transforms source chain events to sdk.Msg
type Transformer interface {
	TransformInboundEvents() error
}
