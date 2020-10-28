package threading

import (
	"github.com/binwen/zero-tools/rescue"
)

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer rescue.Recover()

	fn()
}
