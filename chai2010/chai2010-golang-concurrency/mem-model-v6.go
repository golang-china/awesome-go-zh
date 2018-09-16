// +build ignore

package main

import (
    "sync"
)

func main() {
    var wg sync.WaitGroup
    var limit = make(chan struct{}, 3)
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            limit <- struct{}{} // len(limit) 小于 cap(limit) 才能进入 // HL
            defer func(){ <-limit }() // 退出时 len(limit) 减 1 // HL

            println(id)
        }(i)
    }
    wg.Wait()
}
