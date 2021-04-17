package rosettelog

import (
	"fmt"
)

func Success(msg interface{}) {
	Symbol('+', "2", msg)
}

func Fail(msg interface{}) {
	Symbol('-', "1", msg)
}

func Critical(msg interface{}) {
	Symbol('!', "1", msg)
}

func Symbol(sym rune, color string, msg interface{}) {
	fmt.Println("\u001b[90m[\u001b[3" + color + "m" + string(sym) + "\u001b[90m]\u001b[0m", msg)
}
