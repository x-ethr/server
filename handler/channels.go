package handler

func Channels[Output interface{}]() (output chan *Output, exception chan *Exception, invalid chan *Invalid) {
	output, exception, invalid = make(chan *Output), make(chan *Exception), make(chan *Invalid)

	return
}
