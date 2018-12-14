// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// START // OMIT
const N = 1024 // 纯数字, 无类型 // HL
const str = "this is a 中文 string\n"

var x, y *float
var ch = '\u1234' // HL

/* 定义 T 类型 */
type T struct { a, b int } // HL
var t0 *T = new(T);
t1 := new(T); // 从表达式推导变量的类型 // HL

// 控制结构: 无小括号, 大括号必须 // HL
if len(str) > 0 { ch = str[0] }
