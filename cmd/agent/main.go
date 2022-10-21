package main

import (
	"github.com/jamolpe/kubevisual-agent/internal/api"
)

func main() {
	api := api.New()
	api.Configure()
	api.Listen("1337")
}
