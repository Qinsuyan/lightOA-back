package entity

type FileInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	File []byte `json:"file,omitempty"`
}
