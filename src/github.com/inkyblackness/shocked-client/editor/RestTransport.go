package editor

// RestTransport describes the methods for REST communication.
type RestTransport interface {
	// Get retrieves data from the given URL.
	Get(url string, onSuccess func(jsonString string), onFailure func())
	// Put stores data at the given URL.
	Put(url string, jsonString []byte, onSuccess func(jsonString string), onFailure func())
	// Post requests to add new data at the given URL.
	Post(url string, jsonString []byte, onSucces func(jsonString string), onFailure func())
}
