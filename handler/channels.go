package handler

func Channels[Input interface{}, Body interface{}, Output interface{}]() (input chan *Input, body chan *Body, output chan *Output, exception chan *Exception, invalid chan *Invalid) {
	input, body, output, exception, invalid = make(chan *Input), make(chan *Body, 1), make(chan *Output), make(chan *Exception), make(chan *Invalid)

	return
}
