package input

import "github.com/x-ethr/server/handler/types"

func Channels[Output interface{}]() (output chan *types.Response[Output], exception chan *types.Exception, invalid chan *types.Invalid) {
	output, exception, invalid = make(chan *types.Response[Output]), make(chan *types.Exception), make(chan *types.Invalid)

	return
}
