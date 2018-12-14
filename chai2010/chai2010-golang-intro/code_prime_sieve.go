// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "fmt" // OMIT

// 返回生成自然数序列的管道: 2, 3, 4, ...
func GenerateNatural() chan int { // HL
	ch := make(chan int) // HL
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}() // HL
	return ch
}

// 管道过滤器: 删除能被素数整除的数
func PrimeFilter(in <-chan int, prime int) chan int { // HL
	out := make(chan int) // HL
	go func() {
		for {
			if i := <-in; i%prime != 0 { // HL
				out <- i // HL
			}
		}
	}()
	return out
}

// 素数筛: 菊花链模型
func main() {
	ch := GenerateNatural() // 自然数序列: 2, 3, 4, ... // HL
	for i := 0; i < 10; i++ {
		prime := <-ch // 新出现的素数 // HL
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = PrimeFilter(ch, prime) // 基于新素数构造的过滤器 // HL
	}
}
