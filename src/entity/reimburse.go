package entity

import "time"

type ReimburseInfo struct {
	Id         int                `json:"id"`
	Title      string             `json:"title"`
	Desc       string             `json:"desc"`
	Amount     int                `json:"amount"`
	Auditor    UserInfo           `json:"auditor"`
	Applicant  UserInfo           `json:"applicant"`
	CreatedBy  UserInfo           `json:"createdBy"`
	Status     int                `json:"status"`
	Files      []FileInfo         `json:"files,omitempty"`
	Comments   []ReimburseComment `json:"comments,omitempty"`
	HappenedAt time.Time          `json:"happenedAt"` //格式：2024-06-25T14:34:56+08:00
	UpdatedAt  time.Time          `json:"updatedAt"`
	CreatedAt  time.Time          `json:"createdAt"`
}

type ReimbursePayload struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	Amount     int       `json:"amount"`
	Applicant  int       `json:"applicant"`  //自己作为申请人添加时不传
	HappenedAt time.Time `json:"happenedAt"` //格式：2024-06-25T14:34:56+08:00
}

type ReimburseComment struct {
	Id      int      `json:"id"`
	Comment string   `json:"comment"`
	Creator UserInfo `json:"creator"`
	Time    string   `json:"time"`
}

type ReimburseListFilter struct {
	ListRequest
	Order string `query:"order"` //amount,HappenedAt
	Sort  string `query:"sort"`  //desc,asc

	Title           string    `query:"title"`
	MaxAmount       int       `query:"maxAmount"`
	MinAmount       int       `query:"minAmount"`
	Auditor         int       `query:"auditor"`
	Status          int       `query:"status"`
	Applicant       int       `query:"applicant"`
	HappenedAtStart time.Time `query:"happenedAtStart"`
	HappenedAtEnd   time.Time `query:"happenedAtEnd"`
	CreatedBy       int       `query:"createdBy"` //在查看自己的申请时，这个参数被忽略
}

type ReimburseStatisticFilter struct {
	HappenedAtStart time.Time `query:"happenedAtStart"`
	HappenedAtEnd   time.Time `query:"happenedAtEnd"`
}

type ReimburseAudition struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
	Approve bool   `json:"approve"` //重新提交时，这个参数被忽略
}

type ReimburseFileDeletion struct {
	ReimburseId int   `json:"reimburseId"`
	FileIds     []int `json:"fileIds"`
}

type ReimburseStatistic struct {
	TotalCount     int `json:"totalCount"`
	ApprovedCount  int `json:"approvedCount"`
	RejectedCount  int `json:"rejectedCount"`
	TotalAmount    int `json:"totalAmount"`
	ApprovedAmount int `json:"approvedAmount"`
	RejectedAmount int `json:"rejectedAmount"`
}

// total_requests: 总报销请求数。
// reimbursed_requests: 已报销请求数。
// unreimbursed_requests: 未报销请求数。
// total_amount: 报销总金额。
// daily_totals:
// 每天的总报销金额和报销请求数量。
// department_totals:
// 每个部门的总报销金额和报销请求数量。
type ReimburseStatisticInfo struct {
	ReimburseStatistic
	Daily struct {
		Data []ReimburseStatistic `json:"data"`
		Date []string             `json:"date"`
	} `json:"daily"`
	Department struct {
		Data       []ReimburseStatistic `json:"data"`
		Department []Department         `json:"department"`
	} `json:"department,omitempty"`
}
