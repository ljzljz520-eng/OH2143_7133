package catalog

import (
	"fmt"
	"strings"

	"wedding-sign/internal/model"
)

type GuestMessage struct {
	Greeting string `json:"greeting"`
	Detail   string `json:"detail"`
	Action   string `json:"action"`
	Locale   string `json:"locale"`
}

func BuildGuestMessage(record model.Record, profile model.Profile) GuestMessage {
	locale := profile.Locale
	if locale == "" {
		locale = "en-US"
	}
	if locale == "zh-CN" {
		return GuestMessage{Greeting: "欢迎光临", Detail: fmt.Sprintf("请前往%s", record.Venue), Action: "开始", Locale: locale}
	}
	return GuestMessage{Greeting: "Welcome", Detail: fmt.Sprintf("Please join us at %s", record.Venue), Action: "Begin", Locale: locale}
}

func (m GuestMessage) Complete() bool { return m.Greeting != "" && m.Detail != "" && m.Action != "" }

func (m GuestMessage) Lines() []string { return []string{m.Greeting, m.Detail, m.Action} }

func (m GuestMessage) SearchText() string { return strings.Join(m.Lines(), " ") }

func LocalizeAction(locale string) string {
	if strings.EqualFold(locale, "zh-CN") {
		return "开始"
	}
	if strings.EqualFold(locale, "ja-JP") {
		return "開始"
	}
	return "Begin"
}
