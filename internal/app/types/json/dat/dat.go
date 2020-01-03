package dat

type Object struct {
	ThreadTitle  string    `json:"thread_title"`
	Messages     []Message `json:"messages"`
	LastModified string    `json:"last_modified"`
	Precure      int64     `json:"precure"`
}

type Message struct {
	Num       int    `json:"num"`
	Name      string `json:"name"`
	Mail      string `json:"mail"`
	DateAndId string `json:"date_and_id"`
	Content   string `json:"content"`
}
