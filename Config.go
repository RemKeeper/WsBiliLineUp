package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type 配置文件 struct {
	RoomId           string
	Ps               string
	GuardPrintColor  int
	PsG              string
	GiftPrintColor   int
	PSG              string
	GiftLinePrice    float64
	PsGp             string
	CommonPrintColor int
	PsC              string
	LineKey          string
	PPs              string
}

func GetConfig() (Roomid string, GuardPrintColor, GiftPrintColor int, GiftLinePrice float64, CommonPrintColor int, Linekey string) {
	lineupConfigFile := "./lineupConfig.json"
	Configinfo, OpenErr := os.Open(lineupConfigFile)
	if OpenErr != nil {
		fmt.Println("配置读取错误", OpenErr)
		return
	} else {
		ReadByte := make([]byte, 1024)
		for {
			over, ReadByteErr := Configinfo.Read(ReadByte)
			if over == 0 || ReadByteErr == io.EOF {
				break
			}
			var ConfigSetGet 配置文件
			err := json.Unmarshal(ReadByte[:over], &ConfigSetGet)
			if err != nil {
				return
			}
			return ConfigSetGet.RoomId, ConfigSetGet.GuardPrintColor, ConfigSetGet.GiftPrintColor, ConfigSetGet.GiftLinePrice, ConfigSetGet.CommonPrintColor, ConfigSetGet.LineKey
		}
	}
	return
}

func SetConfig(Roomid string, GuardPrintColor, GiftPrintColor int, GiftLinePrice float64, CommonPrintColor int, Linekey string) bool {
	ConfigJsonType := 配置文件{
		RoomId:           Roomid,
		Ps:               "上面这个是房间号",
		GuardPrintColor:  GuardPrintColor,
		PsG:              "上面这个是舰长队列打印颜色",
		GiftPrintColor:   GiftPrintColor,
		PSG:              "上面这个是礼物队列打印颜色",
		GiftLinePrice:    GiftLinePrice,
		PsGp:             "上面那个是触发礼物排队的单次投喂价格",
		CommonPrintColor: CommonPrintColor,
		PsC:              "上面这个是普通用户打印颜色",
		LineKey:          Linekey,
		PPs:              "上面这个是排队关键词",
	}
	ConfigJson, _ := json.MarshalIndent(ConfigJsonType, "", " ")
	fmt.Println(string(ConfigJson))
	lineupConfig := "./lineupConfig.json"
	_, ReadConfigErr := os.Open(lineupConfig)
	if ReadConfigErr != nil {
		fmt.Println("配置文件不存在，尝试创建")
		_, ConfigErr := os.Create(lineupConfig)
		_, LineCreate := os.Create("./line.json")
		if ConfigErr != nil || LineCreate != nil {
			fmt.Println("配置文件创建失败", ConfigErr.Error(), LineCreate.Error())
			return false
		} else {
			err := os.WriteFile(lineupConfig, ConfigJson, 0666)
			if err != nil {
				fmt.Println("配置文件更新失败", err.Error())
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func SetLine(lp 队列) {
	lineJson, _ := json.MarshalIndent(lp, "", " ")
	lineConfigFile := "./line.json"
	WriteErr := os.WriteFile(lineConfigFile, lineJson, 0666)
	if WriteErr != nil {
		fmt.Println("队列文件更新失败")
	}
}

func GetLine() 队列 {
	lineConfigFile := "./line.json"
	var LineGet 队列
	file, err := os.ReadFile(lineConfigFile)
	if err != nil {
		return 队列{}
	} else {
		_ = json.Unmarshal(file, &LineGet)
		return LineGet
	}
}
