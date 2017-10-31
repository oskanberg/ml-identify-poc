package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/machinebox/sdk-go/tagbox"
	"github.com/machinebox/sdk-go/x/boxutil"
)

func main() {
	var (
		addr = flag.String("addr", ":8080", "address")
	)

	flag.Parse()
	tagbox := tagbox.New("http://mb:8080")

	fmt.Println("Waiting for Facebox to be ready...")
	boxutil.WaitForReady(context.Background(), tagbox)
	fmt.Println("Done!")

	fmt.Println("Go to:", *addr+"...")

	srv := NewServer("./assets", tagbox)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatalln(err)
	}
}
