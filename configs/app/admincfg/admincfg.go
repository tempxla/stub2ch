package admincfg

// 管理者用の設定
// 公開してもまあ問題ないけど、非公開の方がよい
const (
	// 管理者用のセッションクッキーの名前
	LOGIN_COOKIE_NAME = "ADMIN_SESSION_KEY"

	// パスフレーズのダイジェスト
	LOGIN_PASSPHRASE_DIGEST = "e46f2f6a3407dcb62be4d2c6e371cdce914c0d1f0483ab6818327f8518cb7ca8"

	// パスフレーズのリクエストパラメータ名
	LOGIN_PASSPHRASE_PARAM = "ADMIN_PASSPHRASE"

	// シグネチャのリクエストパラメータ名
	LOGIN_SIGNATURE_PARAM = "ADMIN_SIGNATURE"
)
