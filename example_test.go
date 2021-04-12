// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"fmt"
	"net/http"
	"time"

	tingyun3 "github.com/TingYunGo/goagent"
)

func ExampleCreateAction() {
	action, _ := tingyun3.CreateAction("ROUTER", "main.ExampleCreateAction")
	time.Sleep(time.Millisecond * 100)
	action.Finish()
}
func ExampleAction_CreateComponent() {
	action, _ := tingyun3.CreateAction("ROUTER", "main.ExampleAction_CreateComponent")
	component := action.CreateComponent("ExampleAction_CreateComponent")
	subComponent := component.CreateComponent("subcomponent")
	time.Sleep(time.Millisecond * 100)
	subComponent.Finish()
	component.Finish()
	action.Finish()
}
