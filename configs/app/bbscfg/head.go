package bbscfg

var (
	heads = map[string][]byte{
		"new4vip": []byte(`<pre>
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
</pre>`),

		"poverty": []byte(`<pre>
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
</pre>`),
	}
)

func MakeHeadTxt(boardName string) []byte {
	return heads[boardName]
}
