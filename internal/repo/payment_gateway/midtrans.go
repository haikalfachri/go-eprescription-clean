package payment_gateway

import (
	"go-eprescription-clean/pkg/midtrans"
)

// MidtransRepo - Midtrans repository struct.
type MidtransRepo struct {
	*midtrans.SnapClient
}

// NewMidtransRepo - creates a new Midtrans repository.
func NewMidtransRepo(m *midtrans.SnapClient) *MidtransRepo {
	return &MidtransRepo{
		SnapClient: m,
	}
}

func (m *MidtransRepo) CreateSnapTransaction(transactionID string, amount int64, customerName string) (redirectURL, token string, err error) {
	return m.CreateSnap(transactionID, amount, customerName)
}
