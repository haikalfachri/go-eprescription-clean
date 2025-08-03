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
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY") // or from config
	midtrans.Environment = midtrans.Sandbox // or midtrans.Production

	client := snap.Client{}
	client.New(midtrans.ServerKey, midtrans.Sandbox)

	return &SnapClient{Client: client}
}

func (c *SnapClient) CreateSnapTransaction(transactionID string, amount int64, customerName string) (*snap.Response, *midtrans.Error) {
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
		return resp, midErr
	}

	return resp, nil
}
