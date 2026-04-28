package dto

type CreateNoteRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Topic    string   `json:"topic" binding:"required"`
	Tags     []string `json:"tags"`
	FolderID *uint    `json:"folder_id"`
}

type UpdateNoteRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Topic    string   `json:"topic" binding:"required"`
	Tags     []string `json:"tags"`
	FolderID *uint    `json:"folder_id"`
}
