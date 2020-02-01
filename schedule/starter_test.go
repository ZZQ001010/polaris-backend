package main

import (
	"fmt"
	"github.com/Jeffail/tunny"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fn := func(payload interface{}) interface{} {
		fmt.Println(time.Now().Second())
		time.Sleep(time.Duration(2)*time.Second)
		return nil
	}

	//定时任务通用协程池
	pool := tunny.NewFunc(3, fn)
	defer pool.Close()

	go pool.Process(1)
	go pool.Process(2)
	go pool.Process(3)
	go pool.Process(4)
	go func() {
		obj, err := pool.ProcessTimed(5, time.Duration(3)*time.Second)
		fmt.Println(obj, err)
	}()

	time.Sleep(time.Duration(3)*time.Second)
}
