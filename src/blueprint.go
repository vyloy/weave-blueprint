package main

import (
	"github.com/loomnetwork/go-loom/plugin"
	"github.com/loomnetwork/weave-blueprint/src/blueprint"
)

func main() {
	plugin.Serve(blueprint.Contract)
}
