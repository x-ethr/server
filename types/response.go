package types

// Response serves as a data structure for representing an HTTP response.
type Response struct {
	Status  int         // Status represents the HTTP status code of an HTTP response.
	Payload interface{} // Payload is an interface representing the payload data of an HTTP response.
}
