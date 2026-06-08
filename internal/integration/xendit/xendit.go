package xendit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rentalin/internal/errs"
)

func (c *client) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errs.WrapErr("MarshalInvoiceRequest", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/v2/invoices",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, errs.WrapErr("CreateHTTPRequest", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(c.secretKey, "")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, errs.WrapErr("SendInvoiceRequest", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapErr("ReadInvoiceResponse", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"xendit error: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var result CreateInvoiceResponse

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errs.WrapErr("UnmarshalInvoiceResponse", err)
	}

	return &result, nil
}
