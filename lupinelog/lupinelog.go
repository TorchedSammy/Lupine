package lupinelog

import (
	"fmt"

	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})
	L.SetField(mod, "success", luar.New(L, Success))
	L.SetField(mod, "fail", luar.New(L, Fail))
	L.SetField(mod, "critical", luar.New(L, Critical))

	L.Push(mod)

	return 1
}

func Success(msg ...interface{}) {
	Symbol('+', "2", msg...)
}

func Fail(msg ...interface{}) {
	Symbol('-', "1", msg...)
}

func Critical(msg ...interface{}) {
	Symbol('!', "1", msg...)
}

func Symbol(sym rune, color string, msg ...interface{}) {
	fmt.Print("\u001b[90m[\u001b[3" + color + "m" + string(sym) + "\u001b[90m]\u001b[0m")
	fmt.Println(msg...)
}
