package output

import "github.com/x-ethr/server/handler/types"

func Channels[Output interface{}]() (output chan *Output, exception chan *types.Exception) {
	output, exception = make(chan *Output), make(chan *types.Exception)

	return
}
