package server

import "github.com/x-ethr/server/types"

func channels() (output chan *types.Response, exception chan *types.Exception) {
	output, exception = make(chan *types.Response), make(chan *types.Exception)

	return
}
