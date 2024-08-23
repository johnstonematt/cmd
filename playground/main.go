package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/johnstonematt/cmd"
)

func main() {
	if len(os.Args) > 1 {
		server()
	} else {
		client()
	}
}

func server() {
	ttl := time.Second * 10
	timeout := time.After(ttl)

	fmt.Printf("[server] start (running for %s)\n", ttl)

	for {
		n := time.Duration(rand.Int64N(10))
		step := time.After(time.Millisecond * n)

		select {
		case <-step:
			fmt.Printf("\033[2K\r[server] %d", n)
		case <-timeout:
			fmt.Println("\033[2K\r[server] exit")
			return
		}
	}

}

func client() {
	selfExe, err := os.Executable()

	if err != nil {
		panic(err)
	}

	svCmd := []string{selfExe, "server"}

	log.Printf("[client] launching server: %q", svCmd)

	cmd_ := cmd.NewCmdOptions(cmd.Options{Streaming: true}, selfExe, "server")
	//cmd_ := cmd.NewCmdOptions(cmd.Options{Streaming: true}, "echo", "howdy \n partner")

	go func() {
		for output := range cmd_.Stdout {
			fmt.Fprintln(os.Stdout, output)
		}
	}()

	go func() {
		for output := range cmd_.Stderr {
			fmt.Fprintln(os.Stderr, output)
		}
	}()

	log.Println("[client] wait server")

	status := <-cmd_.Start()

	if status.Error != nil {
		log.Fatal(status.Error)
	}

	log.Printf("[client] server status: %v", status)
	log.Println("[client] exit")
}
