package board

type Object struct {
	Subjects []Subject `json:"subjects"`
	Precure  int64     `json:"precure"`
}

type Subject struct {
	ThreadKey    string `json:"thread_key"`
	ThreadTitle  string `json:"thread_title"`
	MessageCount int    `json:"message_count"`
	LastModified string `json:"last_modified"`
}
