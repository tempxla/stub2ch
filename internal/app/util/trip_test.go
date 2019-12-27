package util

import (
	"testing"
)

func TestComputeTrip(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		// 新
		{"ABCDEFGHIJKL", "KqNqEREFnAu9"},  //12
		{"ABCDEFGHIJKLM", "2VMSxHqqU6tf"}, // 13
		// 旧
		{"ABCDEFGHIJK", "aOLjRoi1zs"}, // 11
	}

	for _, tt := range tests {

		actual := ComputeTrip(tt.key)

		if actual != tt.expected {
			t.Errorf("%v : %v : %v", tt.key, tt.expected, actual)
		}
	}
}

func TestComputeTrip12(t *testing.T) {

	tests := []struct {
		key      string
		expected string
	}{
		// 生キー
		{"#9CA39C423D4881A6..", "moussy./hk"},
		{"#ZZA39C423D4881A6..", "???"},       // NOT HEX
		{"#9CA39C423D4881A6H", "LpYKngoFKI"}, // salt 1
		{"#9CA39C423D4881A6", "moussy./hk"},  // salt 0
		// 予約キー
		{"$a", "???"},
		// 新方式
		{"ABCDEFGHIJKL", "KqNqEREFnAu9"},              //12
		{"ABCDEFGHIJKLM", "2VMSxHqqU6tf"},             // 13
		{UTF8toSJISString("あいうえおあ"), "Zca6CIYTvJlI"},  // 12byte
		{UTF8toSJISString("あいうえおあA"), "/33ThS09qMTH"}, // 13byte
		{UTF8toSJISString("あいう朽おあ"), "ZJtW.mI1WEAp"},  // 朽 0x8b 0x80
	}

	for _, tt := range tests {

		actual := computeTrip12(tt.key)

		if actual != tt.expected {
			t.Errorf("%v : %v : %v", tt.key, tt.expected, actual)
		}
	}

}

func TestComputeTripOld(t *testing.T) {

	tests := []struct {
		key      string
		expected string
	}{
		{"", "jPpg5.obl6"},
		{"a", "ZnBI2EKkq."},
		{UTF8toSJISString("あ"), "3zNBOPkseQ"},
		{UTF8toSJISString("あ朽あ"), "YcxY/shYNQ"}, // 朽 0x8b 0x80
	}

	for _, tt := range tests {

		actual := computeTripOld(tt.key)

		if actual != tt.expected {
			t.Errorf("%v : %v : %v", tt.key, tt.expected, actual)
		}
	}
}
