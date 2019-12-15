package config

const (
	// 管理者用のセッションクッキーの名前
	ADMIN_COOKIE_NAME = "ADMIN_SESSION_KEY"

	// パスフレーズのダイジェスト
	ADMIN_PASSPHRASE_DIGEST = "e46f2f6a3407dcb62be4d2c6e371cdce914c0d1f0483ab6818327f8518cb7ca8"

	// パスフレーズのリクエストパラメータ名
	ADMIN_PASSPHRASE_PARAM = "ADMIN_PASSPHRASE"

	// シグネチャのリクエストパラメータ名
	ADMIN_SIGNATURE_PARAM = "ADMIN_SIGNATURE"
)
