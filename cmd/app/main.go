package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/bububa/aliyun-acme-hook/internal/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// rand.Seed(time.Now().UTC().UnixNano())
	ctx := context.Background()
	server := new(cli.App)
	app.NewApp(server)
	if err := server.RunContext(ctx, os.Args); err != nil {
		log.Println(err)
	}
}
