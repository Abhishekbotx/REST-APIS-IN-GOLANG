package main

import ("fmt")

func main() {
	fmt.Println("Hello, World!")
	printer(123)
}

func printer(num int64){
    fmt.Println(num)
}
	// printer(123)