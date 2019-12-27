// 解説
// https://2ch.me/vikipedia/%E3%83%88%E3%83%AA%E3%83%83%E3%83%97
// http://trip.ageruo.jp/trip/triptest.cgi
// https://ja.osdn.net/projects/zerochplus/svn/view/triptest/tags/1.0.1/trip.cgi?view=markup&revision=17&root=zerochplus
// https://ja.osdn.net/projects/naniya/wiki/2chtrip#

package util

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/tempxla/stub2ch/internal/app/util/lib/crypt"
	"regexp"
	"strings"
)

// keyは#より後の文字列
func ComputeTrip(key string) (trip string) {

	if len(key) >= 12 {
		trip = computeTrip12(key)
	} else {
		trip = computeTripOld(key)
	}
	return
}

func computeTrip12(key string) (trip string) {
	mark := key[:1]

	if mark == "#" || mark == "$" {

		re := regexp.MustCompile(`^#([0-9A-Fa-f]{16})([\./0-9A-Za-z]{0,2})$`) // 16進の正規表現 fix
		match := re.FindStringSubmatch(key)

		if match != nil {
			// 生キー
			key2, _ := hex.DecodeString(match[1])
			salt := (match[2] + "..")[:2]

			// 0x80問題再現 :今は実装されていない？
			// x80i := bytes.IndexByte(key2, 0x80)
			// if x80i != -1 {
			// 	key2 = key2[:x80i]
			// }

			t := crypt.Crypt(key2, []byte(salt))
			trip = string(t[len(t)-10:])

		} else {
			// 予約キー
			trip = "???"
		}

	} else {
		// 新仕様
		full := sha1.Sum([]byte(key))
		trip = base64.StdEncoding.EncodeToString(full[:])[:12]
		trip = strings.ReplaceAll(trip, "+", ".")
	}

	return
}

func computeTripOld(key string) (trip string) {

	var salt string
	if key == "" {
		salt = "H." // key == "" の場合は . になってしまうので
	} else {
		salt = (key + "H.")[1:MinInt(len(key)+2, 3)]

		reg := regexp.MustCompile(`[^\.-z]`)
		salt = reg.ReplaceAllString(salt, ".")

		rep := strings.NewReplacer(
			":", "A", ";", "B", "<", "C", "=", "D", ">", "E", "?", "F", "@", "G",
			"[", "a", "\\", "b", "]", "c", "^", "d", "_", "e", "`", "f")
		salt = rep.Replace(salt)
	}

	key2 := []byte(key)
	// 0x80問題再現 :今は実装されていない？
	// x80i := bytes.IndexByte(key2, 0x80)
	// if x80i != -1 {
	// 	key2 = key2[:x80i]
	// }

	t := crypt.Crypt(key2, []byte(salt))
	trip = string(t[len(t)-10:])

	return
}
