package output

import "github.com/x-ethr/server/handler/types"

func Channels[Output interface{}]() (output chan *types.Response[Output], exception chan *types.Exception) {
	output, exception = make(chan *types.Response[Output]), make(chan *types.Exception)

	return
}
