package main

import (
	"fmt"
	"strconv"
)

func KeyBordCtrl(data string, SelectLine *int, Line 队列) {
	switch {
	case data == "Down" || data == "Right":
		*SelectLine++
		PrintLine(SelectLine)
	case data == "Up" || data == "Left":
		*SelectLine--
		PrintLine(SelectLine)
	case data == "Delete":
		DeleteLine(*SelectLine)
	case data == "End":
		*SelectLine = 0
		PrintLine(SelectLine)
	}
}

func PrintLine(SelectLine *int) {
	//fmt.Println("打印", SelectLine, lineup)
	CallClear()
	GuarLineLEN, GiftLineLEN, CommonLineLEN := len(lineup.GuardLine), len(lineup.GiftLine), len(lineup.CommonLine)
	if *SelectLine+1 > GuarLineLEN+GiftLineLEN+CommonLineLEN {
		*SelectLine = 0
	}
	if *SelectLine < (-1) {
		*SelectLine = GuarLineLEN + GiftLineLEN + CommonLineLEN - 1
	}
	//fmt.Println(GuarLineLEN, GiftLineLEN, CommonLineLEN)
	switch {
	case *SelectLine+1 <= GuarLineLEN:
		for k, s := range lineup.GuardLine {
			if *SelectLine == k {
				ColorPrint(s.UserName, 4)
				fmt.Println()
			} else {
				ColorPrint(s.UserName, s.PrintColor)
				fmt.Println()
			}
		}
		for _, s := range lineup.GiftLine {
			ColorPrint(s.UserName+" "+strconv.FormatFloat(s.GiftPrice, 'g', 5, 64)+"元", s.PrintColor)
			fmt.Println()
		}
		for _, s := range lineup.CommonLine {
			ColorPrint(s.UserName, s.PrintColor)
			fmt.Println()
		}
	case GuarLineLEN < *SelectLine+1 && *SelectLine+1 <= GiftLineLEN+GuarLineLEN:
		for _, s := range lineup.GuardLine {
			ColorPrint(s.UserName, s.PrintColor)
			fmt.Println()
		}
		for k, s := range lineup.GiftLine {
			if k == *SelectLine-GuarLineLEN {
				ColorPrint(s.UserName+" "+strconv.FormatFloat(s.GiftPrice, 'g', 5, 64)+"元", 4)
				fmt.Println()
			} else {
				ColorPrint(s.UserName+" "+strconv.FormatFloat(s.GiftPrice, 'g', 5, 64)+"元", s.PrintColor)
				fmt.Println()
			}
		}
		for _, s := range lineup.CommonLine {
			ColorPrint(s.UserName, s.PrintColor)
			fmt.Println()
		}
	case GuarLineLEN+GiftLineLEN < *SelectLine+1 && *SelectLine+1 <= GuarLineLEN+GiftLineLEN+CommonLineLEN:
		for _, s := range lineup.GuardLine {
			ColorPrint(s.UserName, s.PrintColor)
			fmt.Println()
		}
		for _, s := range lineup.GiftLine {
			ColorPrint(s.UserName+" "+strconv.FormatFloat(s.GiftPrice, 'g', 5, 64)+"元", s.PrintColor)
			fmt.Println()
		}
		for k, s := range lineup.CommonLine {
			if k == *SelectLine-GuarLineLEN-GiftLineLEN {
				ColorPrint(s.UserName, 4)
				fmt.Println()
			} else {
				ColorPrint(s.UserName, s.PrintColor)
				fmt.Println()
			}
		}
	}
}
