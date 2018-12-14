package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// unbuffered
	fmt.Fprintf(os.Stdout, "%s, ", "hello")
	// buffered: os.Stdout implements io.Writer
	buf := bufio.NewWriter(os.Stdout) // HL
	// and now so does buf.
	fmt.Fprintf(buf, "%s\n", "world!")
	buf.Flush() // HL
}
