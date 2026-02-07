package sepay

import (
	"context"
	"fmt"
	"net/url"
)

// OrderQueryParams holds optional query parameters for listing orders.
type OrderQueryParams struct {
	PerPage       *int
	Q             *string
	OrderStatus   *string
	CreatedAt     *string
	FromCreatedAt *string
	ToCreatedAt   *string
	CustomerID    *string
	SortCreatedAt *string
}

func (p *OrderQueryParams) toValues() url.Values {
	if p == nil {
		return nil
	}
	v := url.Values{}
	if p.PerPage != nil {
		v.Set("per_page", fmt.Sprintf("%d", *p.PerPage))
	}
	if p.Q != nil {
		v.Set("q", *p.Q)
	}
	if p.OrderStatus != nil {
		v.Set("order_status", *p.OrderStatus)
	}
	if p.CreatedAt != nil {
		v.Set("created_at", *p.CreatedAt)
	}
	if p.FromCreatedAt != nil {
		v.Set("from_created_at", *p.FromCreatedAt)
	}
	if p.ToCreatedAt != nil {
		v.Set("to_created_at", *p.ToCreatedAt)
	}
	if p.CustomerID != nil {
		v.Set("customer_id", *p.CustomerID)
	}
	if p.SortCreatedAt != nil {
		v.Set("sort[created_at]", *p.SortCreatedAt)
	}
	if len(v) == 0 {
		return nil
	}
	return v
}

// OrderService provides access to the order-related endpoints.
type OrderService struct {
	api apiResource
}

// All retrieves a list of orders matching the given query parameters.
func (s *OrderService) All(ctx context.Context, params *OrderQueryParams) (*Response, error) {
	return s.api.doRequest(ctx, "GET", "order", params.toValues(), nil)
}

// Retrieve retrieves the details of a single order by its invoice number.
func (s *OrderService) Retrieve(ctx context.Context, orderInvoiceNumber string) (*Response, error) {
	return s.api.doRequest(ctx, "GET", "order/detail/"+orderInvoiceNumber, nil, nil)
}

// VoidTransaction voids a transaction for the given order invoice number.
func (s *OrderService) VoidTransaction(ctx context.Context, orderInvoiceNumber string) (*Response, error) {
	body := map[string]string{"order_invoice_number": orderInvoiceNumber}
	return s.api.doRequestJSON(ctx, "POST", "order/voidTransaction", body)
}

// Cancel cancels the order with the given invoice number.
func (s *OrderService) Cancel(ctx context.Context, orderInvoiceNumber string) (*Response, error) {
	body := map[string]string{"order_invoice_number": orderInvoiceNumber}
	return s.api.doRequestJSON(ctx, "POST", "order/cancel", body)
}
