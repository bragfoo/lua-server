package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/yuin/gopher-lua"
	"strconv"
)

func main() {
	app := iris.Default()
	app.Post("/lua_script", luaHandler())
	// listen and serve on http://0.0.0.0:8080.
	app.Run(iris.Addr(":8080"))
}

type Pairs struct {
	Key   string `json:"key"`
	Value string `json:"value"` //base64
}

type EvalRequest struct {
	Script  string   `json:"script"`
	NumKeys int      `json:"numkeys"`
	KVs     []Pairs  `json:"kvs"`
	Argv    []string `json:"argv"`
}

type EvalResponse struct {
	KVs []Pairs `json:"kvs"`
}

func luaHandler() iris.Handler {
	return func(ctx iris.Context) {
		
		//run lua vm
		L := lua.NewState()
		defer L.Close()

		//get req
		var req EvalRequest
		if err := ctx.ReadJSON(&req); err != nil {
			panic(err)
		}

		// load file
		if err := L.DoFile("./lua/p50.lua"); err != nil {
			panic(err)
		}

		//init values
		var valus []lua.LValue

		//logic
		kvs := req.KVs
		for _, kv := range kvs {
			vlaueInt, _ := strconv.ParseInt(kv.Value, 10, 64)
			valus = append(valus, lua.LNumber(vlaueInt))
		}

		// run func
		L.CallByParam(lua.P{
			Fn:      L.GetGlobal("p50"),
			NRet:    1,
			Protect: true,
		}, valus...)

		//get lua return
		ret := L.Get(-1)

		//remove value
		L.Pop(1)

		res, _ := ret.(lua.LNumber)
		fmt.Println(int(res))

		//resp
		rsp := &EvalResponse{}
		ctx.JSON(rsp)
	}
}
