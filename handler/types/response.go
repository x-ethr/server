package types

type Response[Output interface{}] struct {
	Code    int
	Payload Output
}
