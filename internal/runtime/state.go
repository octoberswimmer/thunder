package runtime

import "syscall/js"

// currentDiv stores the div element for the current Thunder instance
var currentDiv js.Value

// SetCurrentDiv sets the current div element for the Thunder instance
func SetCurrentDiv(div js.Value) {
	currentDiv = div
}

// GetCurrentDiv returns the current div element for the Thunder instance
func GetCurrentDiv() js.Value {
	return currentDiv
}
