package midtrans

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type SnapClient struct {
	Client snap.Client
}

// New returns a ready-to-use SnapClient
func New() *SnapClient {
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	env := os.Getenv("MIDTRANS_ENVIRONMENT")
	switch env {
	case "production":
		midtrans.Environment = midtrans.Production
	default:
		midtrans.Environment = midtrans.Sandbox
	}

	client := snap.Client{}
	client.New(midtrans.ServerKey, midtrans.Sandbox)

	return &SnapClient{Client: client}
}

func (c *SnapClient) CreateSnap(transactionID string, amount int64, customerName string) (redirectURL, token string, err *midtrans.Error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
		},
	}
	
	resp, midErr := c.Client.CreateTransaction(req)

	if midErr != nil {
		// Log the midtrans error but still return the response if it's usable
		// You can inspect midErr.StatusCode, midErr.Message, etc.
		return "", "", midErr
	}

	return resp.RedirectURL, resp.Token, nil
}


