package xendit

import (
	"context"
	"fmt"
	"os"

	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type XenditClient struct {
	Client *xendit.APIClient
}

// New initializes a new Xendit API client using the secret key.
func New() *XenditClient {
	client := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))

	return &XenditClient{
		Client: client,
	}
}

func (c *XenditClient) CreateInvoice(ctx context.Context, transactionID string, amount int64) (InvoiceURL string, err error) {
	req := *invoice.NewCreateInvoiceRequest(transactionID, float64(amount))

	resp, _, execErr := c.Client.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(req).
		Execute()

	if execErr != nil {
		return "", fmt.Errorf("Error when calling `InvoiceApi.CreateInvoice``: %v\n", execErr.Error())
	}

    return resp.InvoiceUrl, nil
}
