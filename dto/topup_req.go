package dto

type TopUpReq struct {
	Amount int64 `json:"amount"`
	UserID int64 `json:"-"`
}
