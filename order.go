package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	orderNamespace = "/translate/order"
)

// OrderGetResponse defines the response for the OrderGet() endpoint.
type OrderGetResponse struct {
	Order Order `json:"order"`
}

// Order defines an order for the OrderGetResponse.
type Order struct {
	JobsQueued     Int     `json:"jobs_queued"`
	JobsReviewable []Int   `json:"jobs_reviewable"`
	JobsAvailable  []Int   `json:"jobs_available"`
	JobsPending    []Int   `json:"jobs_pending"`
	JobsApproved   []Int   `json:"jobs_approved"`
	JobsRevising   []Int   `json:"jobs_revising"`
	OrderID        Int     `json:"order_id"`
	Credits        Float64 `json:"total_credits"`
	Units          Int     `json:"total_units"`
	Count          Int     `json:"total_jobs"`
	Currency       string  `json:"currency"`
}

// OrderGetRequest defines the request parameters for the OrderGet() endpoint.
type OrderGetRequest struct {
	OrderID int
}

// NewOrderGetRequest creates a new OrderGetRequest with the given id.
func NewOrderGetRequest(orderID int) *OrderGetRequest {
	ogr := &OrderGetRequest{
		OrderID: orderID,
	}
	return ogr
}

// GetOrder retrieves a group of jobs that were previously submitted together by their order id.
func (c *Client) GetOrder(req *OrderGetRequest) (*OrderGetResponse, error) {
	ogr := new(OrderGetResponse)
	err := c.get(orderNamespace+fmt.Sprintf("/%d", req.OrderID), nil, ogr)
	return ogr, err
}

// OrderCancelRequest defines the request parameters for the OrderCancel() endpoint.
type OrderCancelRequest struct {
	OrderID int
}

// NewOrderCancelRequest creates a new OrderCancelRequest with the given id.
func NewOrderCancelRequest(orderID int) *OrderCancelRequest {
	ogr := &OrderCancelRequest{
		OrderID: orderID,
	}
	return ogr
}

// CancelOrder cancels all jobs in an order. Please keep in mind, this endpoint works when all jobs in an order are in available state. This also cancels the order itself.
func (c *Client) CancelOrder(req *OrderCancelRequest) error {
	err := c.delete(orderNamespace+fmt.Sprintf("/%d", req.OrderID), nil, nil)
	return err
}

// OrderCommentsRequest defines the request parameters for the OrderComments() endpoint.
type OrderCommentsRequest struct {
	OrderID int
}

// NewOrderCommentsRequest creates a new OrderCommentsRequest with the given id.
func NewOrderCommentsRequest(orderID int) *OrderCommentsRequest {
	ocr := &OrderCommentsRequest{
		OrderID: orderID,
	}
	return ocr
}

// CommentsResponse defines the response from the OrderComments() and JobComments() endpoints.
type CommentsResponse struct {
	Thread []Comment `json:"thread"`
}

// OrderCommentsResponse defines the response from the OrderComments() endpoint.
type OrderCommentsResponse CommentsResponse

// Comment defines the structure of a Gengo comment response.
type Comment struct {
	Body   string `json:"body"`
	Author string `json:"author,omitempty"`
	Ctime  Time   `json:"ctime,omitempty"`
}

// OrderComments retrieves the comment thread for an order.
func (c *Client) OrderComments(req *OrderCommentsRequest) (*OrderCommentsResponse, error) {
	ocr := new(OrderCommentsResponse)
	err := c.get(orderNamespace+fmt.Sprintf("/%d/comments", req.OrderID), nil, ocr)
	return ocr, err
}

// AddOrderCommentRequest defines the request parameters for the AddOrderComment() endpoint.
type AddOrderCommentRequest struct {
	OrderID int    `json:"id"`
	Body    string `json:"body"`
}

// NewAddOrderCommentRequest creates a new AddOrderCommentRequest with the given id.
func NewAddOrderCommentRequest(orderID int, comment string) *AddOrderCommentRequest {
	r := &AddOrderCommentRequest{
		OrderID: orderID,
		Body:    comment,
	}
	return r
}

// AddOrderComment submits a new comment to the order’s comment thread.
func (c *Client) AddOrderComment(req *AddOrderCommentRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = c.post(orderNamespace+fmt.Sprintf("/%d/comment", req.OrderID), bytes.NewReader(b), nil)
	return err
}
