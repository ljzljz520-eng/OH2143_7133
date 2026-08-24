package catalog

import "strings"

type Labels struct {
	Start      string
	Fullscreen string
	Close      string
	Venue      string
}

func LabelsFor(locale string) Labels {
	switch strings.ToLower(strings.TrimSpace(locale)) {
	case "zh-cn":
		return Labels{Start: "开始", Fullscreen: "全屏", Close: "关闭", Venue: "宴会厅"}
	case "ja-jp":
		return Labels{Start: "開始", Fullscreen: "全画面", Close: "閉じる", Venue: "会場"}
	default:
		return Labels{Start: "Begin", Fullscreen: "Fullscreen", Close: "Close", Venue: "Venue"}
	}
}

func (l Labels) Complete() bool {
	return l.Start != "" && l.Fullscreen != "" && l.Close != "" && l.Venue != ""
}

func (l Labels) All() []string { return []string{l.Start, l.Fullscreen, l.Close, l.Venue} }
