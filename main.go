package main

import (
	"github.com/kataras/iris"
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
		var req EvalRequest
		if err := ctx.ReadJSON(&req); err != nil {
			panic(err)
		}

		rsp := &EvalResponse{}
		ctx.JSON(rsp)
	}
}
