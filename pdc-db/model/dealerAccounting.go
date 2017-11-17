package model

// DealerAccounting Result
type DealerAccounting struct {
	DealerID    uint
	ExternalID  string
	ArticleID   uint
	ArticleText string
	Price       float32
	Currency    string
}
