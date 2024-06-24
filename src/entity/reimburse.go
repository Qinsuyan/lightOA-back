package entity

type ReimburseInfo struct {
	Id        int                `json:"id"`
	Title     string             `json:"title"`
	Desc      string             `json:"desc"`
	Amount    int                `json:"amount"`
	Auditor   UserInfo           `json:"auditor"`
	Applicant UserInfo           `json:"applicant"`
	Status    int                `json:"status"`
	Files     []FileInfo         `json:"files"`
	Comments  []ReimburseComment `json:"comments"`
}

type ReimbursePayload struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Amount int    `json:"amount"`
	Files  []int  `json:"files"`
}

type ReimburseComment struct {
	Id      int      `json:"id"`
	Comment string   `json:"comment"`
	Creator UserInfo `json:"creator"`
	Time    string   `json:"time"`
}
