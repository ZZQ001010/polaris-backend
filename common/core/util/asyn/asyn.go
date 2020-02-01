package asyn

import "fmt"

func Execute(fn func()){
	go func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}

		fn()
	}()
}