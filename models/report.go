package models

type BestSellingProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Count       int    `json:"count"`
}

type SalesSummary struct {
	TotalRevenue       int                 `json:"total_revenue"`
	TotalTransaction   int                 `json:"total_transaction"`
	BestSellingProduct *BestSellingProduct `json:"best_selling_product"`
}
