package main

type Room struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		RoomID     int   `json:"room_id"`
		ShortID    int   `json:"short_id"`
		UID        int   `json:"uid"`
		IsHidden   bool  `json:"is_hidden"`
		IsLocked   bool  `json:"is_locked"`
		IsPortrait bool  `json:"is_portrait"`
		LiveStatus int   `json:"live_status"`
		Encrypted  bool  `json:"encrypted"`
		LiveTime   int64 `json:"live_time"`
	} `json:"data"`
}

type 队列 struct {
	GuardLine  []*Line
	GiftLine   []*GiftLine
	CommonLine []*Line
}

type Line struct {
	Uid        int
	UserName   string
	PrintColor int
}

type GiftLine struct {
	Uid        int
	UserName   string
	PrintColor int
	GiftPrice  float64
}
