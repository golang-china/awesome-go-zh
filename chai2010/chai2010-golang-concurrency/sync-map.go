// +build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
    var m sync.Map // HL
    for i := 0; i < 3; i++ {
        go func(i int) {
            for j := 0; ; j++ {
                m.Store(i, j) // HL
            }
        }(i)
    }
    for i := 0; i < 10; i++ {
        m.Range(func(key, value interface{}) bool { // HL
            fmt.Printf("%d: %d\t", key, value)
            return true // HL
        })
        time.Sleep(time.Second)
		fmt.Println()
    }
}
