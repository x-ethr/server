package input

import "github.com/x-ethr/server/handler/types"

func channels() (output chan *types.Response, exception chan *types.Exception, invalid chan *types.Invalid) {
	output, exception, invalid = make(chan *types.Response), make(chan *types.Exception), make(chan *types.Invalid)

	return
}
