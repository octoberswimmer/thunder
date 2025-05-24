package components

import (
	"github.com/octoberswimmer/masc"
	"github.com/octoberswimmer/masc/elem"
)

// Container renders a simple div container for layout purposes.
// This provides a basic wrapper to avoid direct elem.Div usage in applications.
func Container(children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	if len(children) == 0 {
		return nil
	}
	return elem.Div(children...)
}

// ContainerWithClass renders a div container with custom CSS classes.
func ContainerWithClass(class string, children ...masc.MarkupOrChild) masc.ComponentOrHTML {
	if len(children) == 0 {
		return nil
	}
	args := []masc.MarkupOrChild{
		masc.Markup(masc.Class(class)),
	}
	args = append(args, children...)
	return elem.Div(args...)
}
