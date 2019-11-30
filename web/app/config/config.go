package config

const (
	PROJECT_ID = "stub2ch"

	DAT_DATE_LAYOUT = "2006/01/02"
	DAT_TIME_LAYOUT = "15:04:05.000"
	// 名前<>メール欄<>年/月/日(曜) 時:分:秒.ミリ秒 ID:hogehoge0<> 本文 <>スレタイ
	DAT_FORMAT = "%s<>%s<>%s(%s) %s ID:%s<> %s <>%s"
)

var (
	WEEK_DAYS_JP = [...]string{"日", "月", "火", "水", "木", "金", "土"}
)
