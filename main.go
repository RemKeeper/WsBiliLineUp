package main

import (
	"encoding/json"
	"fmt"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var lineup 队列

func main() {
	NaNSelectLine := 0
	var SelectLine int
	KeyBordString := make(chan string, 1)

	//IsOnlyGift := 0

	RoomId, GuardPrintColor, GiftPrintColor, GiftLinePrice, CommonPrintColor, Linekey, IsOnlyGift := GetConfig()

	if len(RoomId) < 1 || GuardPrintColor < 1 || GiftPrintColor < 1 || GiftLinePrice < 0.01 || CommonPrintColor < 1 || len(Linekey) < 1 {
		RoomId = "2233"
		fmt.Println("请输入bilibili房间号并回车")
		_, _ = fmt.Scanln(&RoomId)
		if RoomId == "11365" {
			fmt.Println("拒绝服务")
			os.Exit(444)
		}
		fmt.Println("排队关键词，默认为”排队“，不使用自定义请直接回车，多个关键词请使用任意符号分隔")
		_, _ = fmt.Scanln(&Linekey)
		if len(Linekey) == 0 {
			Linekey = "排队"
		}
		fmt.Println("是否开启仅限礼物及舰长排队（开启后只有礼物用户及舰长会加入队列）输入1启用")
		_, _ = fmt.Scanln(&IsOnlyGift)
		fmt.Println("自定义颜色编号")
		for i := 1; i < 7; i++ {
			ColorPrint("这是测试字符"+" "+"编号"+strconv.Itoa(i), i)
		}
		for i := 9; i < 15; i++ {
			ColorPrint("这是测试字符"+" "+"编号"+strconv.Itoa(i), i)
		}
		fmt.Println()

		fmt.Print("是否自定义颜色？ 输入1开始自定义 不输入将采用默认[ ")
		ColorPrint("舰长", 27)
		ColorPrint("礼物", 15)
		ColorPrint("普通", 10)
		fmt.Println(" ]配色")
		var EnColorSet int
		_, _ = fmt.Scanln(&EnColorSet)
		if EnColorSet != 1 {
			GuardPrintColor = 27
			GiftPrintColor = 15
			CommonPrintColor = 10
		} else {
			fmt.Println("请输入舰长队列颜色[编号]")
			_, _ = fmt.Scanln(&GuardPrintColor)
			fmt.Println("请输入礼物队列颜色[编号]")
			_, _ = fmt.Scanln(&GiftPrintColor)
			fmt.Println("请输入普通队列颜色[编号]")
			_, _ = fmt.Scanln(&CommonPrintColor)
		}
		fmt.Println("请输入礼物队列触发价格(RMB)，只有单次或累计礼物大于此价格才会被加入")
		_, _ = fmt.Scanln(&GiftLinePrice)
		if SetConfig(RoomId, GuardPrintColor, GiftPrintColor, GiftLinePrice, CommonPrintColor, Linekey, IsOnlyGift) {
			fmt.Println("配置信息已保存，下次启动将自动读取")
		} else {
			fmt.Println("配置文件保存失败")
		}

	}

	//房间信息解析

	RoomInfo := GetRealRoomInfo(RoomId)
	switch {
	case RoomInfo.Code != 0:
		fmt.Println("房间号错误", RoomInfo.Msg, RoomInfo.Message)
		os.Exit(0)
	case RoomInfo.Data.IsHidden != false:
		fmt.Println("隐藏的房间")
		os.Exit(0)
	case RoomInfo.Data.IsLocked != false:
		fmt.Println("被封禁或锁定的房间")
		os.Exit(0)
	case RoomInfo.Data.IsPortrait != false || RoomInfo.Data.Encrypted != false:
		fmt.Println("未知的房间状态，输入1尝试解析弹幕（默认不解析）")
		UserSelect := 0
		_, _ = fmt.Scanln(&UserSelect)
		if UserSelect != 1 {
			os.Exit(0)
		} else {
			break
		}
	default:
		fmt.Println("房间真实ID为", RoomInfo.Data.RoomID)
		fmt.Println("用户uid为", RoomInfo.Data.UID)
		if RoomInfo.Data.LiveStatus == 1 {
			fmt.Println("用户直播状态:直播中")
			fmt.Println("开播时间", time.Unix(RoomInfo.Data.LiveTime, 0))
			fmt.Println("已直播", time.Now().Sub(time.Unix(RoomInfo.Data.LiveTime, 0)))
		} else {
			fmt.Println("用户直播状态:未直播")
		}
		fmt.Println("基础状态检查完成，开始解析弹幕")
		LineTemp := GetLine()
		if len(LineTemp.CommonLine) != 0 || len(LineTemp.GuardLine) != 0 || len(LineTemp.GiftLine) != 0 {
			lineSelect := ""
			fmt.Println("检测到上一次队列缓存，直接回车启用，不启用请输入”D“并回车将删除文件")
			_, _ = fmt.Scanln(&lineSelect)
			if lineSelect == "D" || lineSelect == "d" {
				err := os.Rename("./line.json", "lineback.json")
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("已删除上次队列，即将开始弹幕解析")
				time.Sleep(time.Second * 2)
				CallClear()
			} else {
				lineup = LineTemp
				CallClear()
				for _, s := range lineup.GuardLine {
					ColorPrint(s.UserName, s.PrintColor)
					fmt.Println()
				}
				for _, s := range lineup.GiftLine {
					ColorPrint(s.UserName, s.PrintColor)
					fmt.Println()
				}
				for _, s := range lineup.CommonLine {
					ColorPrint(s.UserName, s.PrintColor)
					fmt.Println()
				}
			}
		}
	}

	//弹幕解析部分

	Newbarrage := client.NewClient(strconv.Itoa(RoomInfo.Data.RoomID))
	Newbarrage.OnDanmaku(func(danmaku *message.Danmaku) {
		if strings.Contains(Linekey, danmaku.Content) {
			if danmaku.Sender.GuardLevel >= 3 {
				lineup.GuardLine = append(lineup.GuardLine, &Line{
					Uid:        danmaku.Sender.Uid,
					UserName:   danmaku.Sender.Uname,
					PrintColor: GuardPrintColor,
				})
				lineup.GuardLine = removeRepeatElement(lineup.GuardLine)
			} else {
				if IsOnlyGift != 1 {
					lineup.CommonLine = append(lineup.CommonLine, &Line{
						Uid:        danmaku.Sender.Uid,
						UserName:   danmaku.Sender.Uname,
						PrintColor: CommonPrintColor,
					})
					lineup.CommonLine = removeRepeatElement(lineup.CommonLine)
				}
			}
			SetLine(lineup)

			PrintLine(&NaNSelectLine)
		}

		if danmaku.Content == "取消排队" {
			for lineCount, User := range lineup.GuardLine {
				if User.UserName == danmaku.Sender.Uname {
					lineup.GuardLine = append(lineup.GuardLine[:lineCount], lineup.GuardLine[lineCount+1:]...)
				}
			}
			for lineCount, User := range lineup.GiftLine {
				if User.UserName == danmaku.Sender.Uname {
					lineup.GiftLine = append(lineup.GiftLine[:lineCount], lineup.GiftLine[lineCount+1:]...)
				}
			}
			for lineCount, User := range lineup.CommonLine {
				if User.UserName == danmaku.Sender.Uname {
					lineup.CommonLine = append(lineup.CommonLine[:lineCount], lineup.CommonLine[lineCount+1:]...)
				}
			}
			SetLine(lineup)
			PrintLine(&NaNSelectLine)
		}
	})

	var GiftLineTemp []GiftLine
	Newbarrage.OnGift(func(gift *message.Gift) {
		if gift.GuardLevel >= 3 {
			if float64((gift.Num*gift.Price)/1000) > GiftLinePrice {
				lineup.GuardLine = append(lineup.GuardLine, &Line{
					Uid:        gift.Uid,
					UserName:   gift.Uname,
					PrintColor: GuardPrintColor,
				})
			}
		} else {
			if GfitIsExist(gift.Uid, lineup.GiftLine) {
				for _, line := range lineup.GiftLine {
					if gift.Uid == line.Uid {
						line.GiftPrice += float64((gift.Num * gift.Price) / 1000)
					}
				}
			} else {
				GiftLineTemp = append(GiftLineTemp, GiftLine{
					Uid:        gift.Uid,
					UserName:   gift.Uname,
					PrintColor: GiftPrintColor,
					GiftPrice:  float64((gift.Num * gift.Price) / 1000),
				})
			}
		}
		for GftempKey, line := range GiftLineTemp {
			if line.GiftPrice >= GiftLinePrice {
				lineup.GiftLine = append(lineup.GiftLine, &line)
				GiftLineTemp = append(GiftLineTemp[:GftempKey], GiftLineTemp[GftempKey+1:]...)
			}
		}
		lineup.GiftLine = removeRepeatGifElement(lineup.GiftLine)
		SetLine(lineup)
		PrintLine(&NaNSelectLine)
	})

	Newbarrage.OnGuardBuy(func(buy *message.GuardBuy) {
		lineup.GuardLine = append(lineup.GuardLine, &Line{
			Uid:        buy.Uid,
			UserName:   buy.Username,
			PrintColor: GuardPrintColor,
		})
		lineup.GuardLine = removeRepeatElement(lineup.GuardLine)
		SetLine(lineup)
		PrintLine(&NaNSelectLine)
	})
	Newbarrage.Start()
	go KeyBordHook(KeyBordString)

	go func() {
		for {
			data, _ := <-KeyBordString
			KeyBordCtrl(data, &SelectLine, lineup)
		}
	}()

	WebServerTwo()
}

func GetRealRoomInfo(RoomId string) Room {
	RoomUrl := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + RoomId
	RealRoomId, err := http.Get(RoomUrl)
	if err != nil {
		fmt.Println("请求房间号错误，请检查您的网络环境", err.Error())
	} else {
		IdBody, _ := io.ReadAll(RealRoomId.Body)
		var RoomConfig Room
		_ = json.Unmarshal(IdBody, &RoomConfig)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(RealRoomId.Body)
		return RoomConfig
	}
	return Room{}
}

func GfitIsExist(UserUid int, AllLine []*GiftLine) bool {
	for _, line := range AllLine {
		if UserUid == line.Uid {
			return true
		}
	}
	return false
}

func DeleteLine(SelectLine int) {
	GuarLineLEN, GiftLineLEN, CommonLineLEN := len(lineup.GuardLine), len(lineup.GiftLine), len(lineup.CommonLine)
	switch {
	case SelectLine+1 <= GuarLineLEN:
		for k := range lineup.GuardLine {
			if SelectLine == k {
				lineup.GuardLine = append(lineup.GuardLine[:k], lineup.GuardLine[k+1:]...)
				SetLine(lineup)
				PrintLine(&SelectLine)

			}
		}
	case GuarLineLEN < SelectLine+1 && SelectLine+1 <= GiftLineLEN+GuarLineLEN:
		for k := range lineup.GiftLine {
			if k == SelectLine-GuarLineLEN {
				lineup.GiftLine = append(lineup.GiftLine[:SelectLine-GuarLineLEN], lineup.GiftLine[SelectLine-GuarLineLEN+1:]...)
				SetLine(lineup)
				PrintLine(&SelectLine)

			}
		}
	case GuarLineLEN+GiftLineLEN < SelectLine+1 && SelectLine+1 <= GuarLineLEN+GiftLineLEN+CommonLineLEN:
		for k := range lineup.CommonLine {
			if k == SelectLine-GuarLineLEN-GiftLineLEN {
				lineup.CommonLine = append(lineup.CommonLine[:SelectLine-GuarLineLEN-GiftLineLEN], lineup.CommonLine[SelectLine-GuarLineLEN-GiftLineLEN+1:]...)
				SetLine(lineup)
				PrintLine(&SelectLine)
			}
		}
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	var FontString string
	dir, err := os.ReadDir("./")
	if err != nil {
		return
	}
	for _, info := range dir {
		if strings.Contains(info.Name(), ".ttf") {
			FontString = info.Name()
		}
	}

	htmlOne := "<head><meta charset=\"utf-8\"><style>@font-face {font-family:name;src: local('./Honkai-zh-cn.ttf'), url('http://127.0.0.1:100/font/" + FontString + "') format('woff');sRules}*{padding: 0px;margin: 0px;}li{list-style:"
	listyle := "\"none\""
	htmlTwo := ";font-family:name;}</style></head><body><ol id=\"father\" style=\"font-size: 40px;\">"
	var lihtml string
	for _, s := range lineup.GuardLine {
		lihtml += "<li style=\"color: " + liColor(s.PrintColor) + ";\">" + s.UserName + "</li>"
	}
	for _, s := range lineup.GiftLine {
		lihtml += "<li style=\"color: " + liColor(s.PrintColor) + ";\">" + s.UserName + " " + strconv.FormatFloat(s.GiftPrice, 'g', 5, 64) + "元" + "</li>"
	}
	for _, s := range lineup.CommonLine {
		lihtml += "<li style=\"color: " + liColor(s.PrintColor) + ";\">" + s.UserName + "</li>"
	}
	htmlThree := "</ol></body><script>function myrefresh(){window.location.reload();};setTimeout('myrefresh()',1000);</script>"
	html := htmlOne + listyle + htmlTwo + lihtml + htmlThree
	fmt.Fprint(w, html)
}

func liColor(ColorInt int) string {
	switch {
	case 0 < ColorInt && ColorInt < 7:
		return []string{"#0644ff", "#12920e", "#3a96dd", "#ff1a2d", "#bb1fd3", "#c19c00"}[ColorInt-1]
	case 8 < ColorInt && ColorInt < 15:
		return []string{"#3b78ff", "#16c60c", "#64dddd", "#e74856", "#b4009e", "#f9f1a5"}[ColorInt-9]
	default:
		return "aqua"
	}
}
