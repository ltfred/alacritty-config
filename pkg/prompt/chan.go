package prompt

var Quit chan bool

func init() {
	Quit = make(chan bool, 1)
}
