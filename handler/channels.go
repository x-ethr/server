package handler

func Channels[Body interface{}, Output interface{}]() (body chan *Body, output chan *Output, exception chan *Exception, invalid chan *Invalid) {
	body, output, exception, invalid = make(chan *Body, 1), make(chan *Output), make(chan *Exception), make(chan *Invalid)

	return
}
