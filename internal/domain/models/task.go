package models

type Task struct {
	ID              string   `bson:"_id,omitempty"`
	Logins          []string `bson:"logins"`
	ApprovalTokens  []string `bson:"approval_tokens"`
	Title           string   `bson:"title"`
	Description     string   `bson:"description"`
	InitiatorLogin  string   `bson:"initiator_login"`
	CurrApprovalNum int      `bson:"curr_approval_num"`
	Status          int      `bson:"status"`
}
