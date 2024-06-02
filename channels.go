package server

import "github.com/x-ethr/server/types"

func channels() (output chan *types.Response, redirect chan *types.Redirect, exception chan *types.Exception) {
	output, redirect, exception = make(chan *types.Response), make(chan *types.Redirect), make(chan *types.Exception)

	return
}
