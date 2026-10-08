package bank

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type bankClient struct {
	bankURL 		string
	httpClient 		*http.Client
}

func NewBankClient(bankURL string, httpClient *http.Client) (*bankClient) {
	return &bankClient{
		bankURL: bankURL,
		httpClient: httpClient,
	}
}

var ErrRequestBank = errors.New("failed to create payment session")

type PaymentRequest struct{
	OrderID 	int64 		`json:"order_id"`
	Amount 		int64 		`json:"amount"`
}

type PaymentResponse struct {
	PaymentID 	int64 		`json:"payment_id"`
	PaymentURL 	string 		`json:"payment_url"`
}

func (b *bankClient) CreatePaymentURL(orderID int64, amount int64) (string, error) {
	request := PaymentRequest{
		OrderID: orderID,
		Amount: amount,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRequestBank, err)
	}

	req, err := http.NewRequest(http.MethodPost, b.bankURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRequestBank, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRequestBank, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w: bank call returned status %d", ErrRequestBank, resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrRequestBank, err)
	}

	var response PaymentResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("%w: %w", ErrRequestBank, err)
	}

	return response.PaymentURL, nil

}