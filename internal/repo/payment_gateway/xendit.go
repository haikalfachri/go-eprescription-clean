package payment_gateway

import (
	"context"

	"go-eprescription-clean/pkg/xendit"
)

// XenditRepo - Xendit repository struct.
type XenditRepo struct {
	*xendit.XenditClient
}

// NewXenditRepo - creates a new Xendit repository.
func NewXenditRepo(x *xendit.XenditClient) *XenditRepo {
	return &XenditRepo{
		XenditClient: x,
	}
}

func (x *XenditRepo) CreateInvoiceTransaction(ctx context.Context, transactionID string, amount int64) (string, error) {
	return x.CreateInvoice(ctx, transactionID, amount)
}

