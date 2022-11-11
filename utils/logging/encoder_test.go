package logging

import (
	"fmt"
	"runtime"
	"testing"
)

func FuncA() string {
	pc, _, _, _ := runtime.Caller(0)
	FuncAName := FuncName(pc)
	ppc, _, _, _ := runtime.Caller(1)
	PFuncAName := FuncName(ppc)
	callInfo := fmt.Sprint(PFuncAName, " called ", FuncAName)
	return callInfo
}

func TestFuncName(t *testing.T) {
	expect := "TestFuncName called FuncA"
	r := FuncA()
	if r != expect {
		t.Error(r)
	}
}
