package head

var (
	head_new4vip = []byte(`<pre>
　　／⌒ヽ
　 ∩ ^ω^) な ん だ
　 |　 ⊂ﾉ
　 |　＿_⊃
　 し′

　 ／⌒ヽ
　(^ω^ ∩　う そ か
　 (⊃　 |
　⊂＿_　|
　　 　` + "`" + `Ｊ

　　 ／⌒ヽ
　　(　　　) おっおっ
　 ／　　_つ　おっ
　(_(_⌒)′
　 ∪(ノ
</pre>`)

	head_poverty = []byte(`<pre>
　　　　　　　　＼　　ヽ　　　　　! |　　　　 /
　　　　　＼　　　　ヽ　　　ヽ　　　　　　　/　　　　/　　 　 　 ／
　　　　　　んああぁぁああぁああああぁぁぁああああ！！！！！
　　　　　　　　＼　　　　　　　　　　｜　 　 　 　 /　　　／
　　　　　　　　　 　 　 　 　 　 　 　 ,ｲ
￣　--　　=　＿　　　　　　　　 　 / |　　　　　　　　　　　　　 --'''''''
　　　　　　　　　　,,, 　 　 ,r‐､λノ　 ﾞi､_,､ﾉゝ　　　　　-　￣
　　　　　　　　　　　　　　ﾞl　　 　 　　 　 　 ﾞ､_
　　　　　　　　　　　　　 .j´　.　.／⌒ヽ　　　（.
　　　　─　　　＿　　─ {　 　 (´ん` + "`" + `#）　　 /─　　　＿　　　　　─
　　　　　　　　　　　　　　 ).　 c/　　 ,つ 　 ,l~
　　　　　　　　　　　　　 ´y　　｛ ,、 ｛　 　 <
　　　　　　　　　　　　　　 ゝ 　 lﾉ ヽ,)　　 ,
</pre>`)
)

func MakeHeadTxt(boardName string) []byte {
	switch boardName {
	case "news4vip":
		return head_new4vip
	case "poverty":
		return head_poverty
	default:
		return nil
	}

	return nil
}
