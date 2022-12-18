package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func ColorPrint(s string, i int) { //设置终端字体颜色
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	fmt.Print(s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, err := CloseHandle.Call(handle)
	if err != nil {
		return
	}
}

func removeRepeatElement(list []*Line) []*Line {
	// 创建一个临时map用来存储数组元素
	temp := make(map[int]bool)
	index := 0
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		_, ok := temp[v.Uid]
		if ok {
			list = append(list[:index], list[index+1:]...)
		} else {
			temp[v.Uid] = true
		}
		index++
	}
	return list
}
func removeRepeatGifElement(list []*GiftLine) []*GiftLine {
	// 创建一个临时map用来存储数组元素
	temp := make(map[int]bool)
	index := 0
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		_, ok := temp[v.Uid]
		if ok {
			list = append(list[:index], list[index+1:]...)
		} else {
			temp[v.Uid] = true
		}
		index++
	}
	return list
}

var clear map[string]func() //create a map for storing clear funcs

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
}
