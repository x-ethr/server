package handler

import "github.com/x-ethr/server/handler/types"

func channels() (output chan *types.Response, exception chan *types.Exception) {
	output, exception = make(chan *types.Response), make(chan *types.Exception)

	return
}
