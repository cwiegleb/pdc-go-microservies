package model

// DealerAccounting Result
type DealerAccounting struct {
	DealerID   uint
	ExternalID string
	ArticleID  uint
	Price      float32
	Currency   string
}
