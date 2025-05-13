package model

type ResumeEvent struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Attachment string `json:"attachment,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	Status     string `json:"status,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

func (r *ResumeEvent) GetKey() string {
	return r.UserID
}

func (r *ResumeEvent) GetId() int {
	return r.ID
}
