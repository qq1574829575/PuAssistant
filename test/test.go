package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Pu后台系统已启动，使用help查看所有指令")
	var command, p1, p2, p3, p4, p5 string
	var taskChs []chan int //所有任务协程的管道数组

Loop:
	for {
		_, _ = fmt.Scanf("%s %s %s %s %s %s", &command, &p1, &p2, &p3, &p4, &p5) //%v按原数据格式 %s输入字符型 %d整数型

		switch command {
		case "exit":
			break Loop
		case "help":
			go help()
			break
		case "newTask":
			ch := make(chan int)
			go Run(ch)
			taskChs = append(taskChs, ch)
			break
		case "exitTask":
			index, err := strconv.Atoi(p1) //传入的第二个参数为要退出任务管道的下标
			if err == nil {
				if len(taskChs) == 0 {
					fmt.Println("当前没有任务在运行!")
					break
				}
				if index < 0 || index > len(taskChs)-1 {
					//如果输入的任务管道下标不在任务数组下标范围内则报错
					fmt.Println("任务不存在!")
					break
				}
				ch := taskChs[index] //获取到指定任务管道
				ch <- 0
				if <-ch == 1 {
					fmt.Printf("退出任务成功，管道数组下标: %v,%v\n", index, taskChs[index])
					taskChs = append(taskChs[:index], taskChs[(index+1):]...) //删除index处的任务管道
					close(ch)                                                 //关闭管道
				} else {
					fmt.Printf("退出任务失败，管道数组下标: %v,%v\n", index, taskChs[index])
				}
			} else {
				fmt.Println("请传入int参数!")
			}
			break
		case "showTaskChs":
			fmt.Println(taskChs)
			break
		case "showTaskStatus":
			//ch <- 1
			//fmt.Println("Run send msg with code",<-ch)
			break
		case "":
			break
		default:
			fmt.Println(command, p1, p2, p3, p4, p5)
		}
		//重置command和所有参数
		command = ""
		p1 = ""
		p2 = ""
		p3 = ""
		p4 = ""
		p5 = ""
	}

}
func help() {
	fmt.Println("1.exit:退出所有正在运行中的协程任务后，退出系统")
}

func Run(done chan int) {
	count := 0
Loop:
	for {
		select {
		case msg := <-done:
			if msg == 1 {
				fmt.Println("Run receive msg code:", msg)
				done <- 1
			}
			if msg == 0 {
				fmt.Println("Run receive msg code:", msg)
				fmt.Println("Run is exiting")
				done <- 1
				break Loop
			}
			break
		default:
		}

		time.Sleep(time.Second * 1)
		count++
		if count == 5 {
			fmt.Println("do something")
			count = 0
		}
	}
}

func test(ch <-chan int) {
	for x := range ch {
		switch x {
		case 1:
			fmt.Printf("此任务正在运行且运行正常")
			break
		}
	}
}
