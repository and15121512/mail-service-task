package models

type Task struct {
	ID              string
	Logins          []string
	ApprovalTokens  []string
	Title           string
	Description     string
	InitiatorLogin  string
	CurrApprovalNum int
	Status          int
}
