package notes

type NoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
