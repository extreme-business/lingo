package token

type Token struct {
	EmailHash       []byte `json:"email_hash"`
	CreateTimestamp int64  `json:"create_timestamp"`
}
