package util

import (
	"time"
)

// 使用说明：
// someFunc：运行的函数，milliseconds：间隔时间毫秒值，async：是否异步运行
// 使用一个变量来接收返回的chan，当往chan发true时，定时器会停止
func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {
	// 间隔
	interval := time.Duration(milliseconds) * time.Millisecond
	// 定时器，到时间时，会有个time信号产生到ticker.C
	ticker := time.NewTicker(interval)
	// 通道，接收外部信号来控制定时器的停止
	clear := make(chan bool)

	//（多线程）
	go func() {
		// 死循环
		for {
			select {
			case <-ticker.C: // 时间到时运行
				if async {
					go someFunc() // 异步运行（多线程）
				} else {
					someFunc()
				}
			case <-clear: // 有信号时，停止定时器
				ticker.Stop()
				return
			}
		}
	}()
	// 返回通道，用来接收控制
	return clear
}

//func main() {
//
//	// A counter for the number of times we print
//	printed := 0
//
//	// We call set interval to print Hello World forever
//	// every 1 second
//	interval := setInterval(func() {
//		fmt.Println("Hello World")
//		printed++
//	}, 1000, false)
//
//	// If we wanted to we had a long running task (i.e. network call)
//	// we could pass in true as the last argument to run the function
//	// as a goroutine
//
//	// Some artificial work here to wait till we've printed
//	// 5 times
//	for {
//		if printed == 5 {
//			// Stop the ticket, ending the interval go routine
//			interval <- true
//			return
//		}
//	}
//
//}
