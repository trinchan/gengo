package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	jobsNamespace = "/translate/jobs"
)

type PostJobsRequest struct {
	Jobs         []*JobRequest `json:"jobs"`
	GroupComment *string       `json:"comment,omitempty"`
}

type PostJobsRequestOption func(*PostJobsRequest)

func WithGroupComment(s string) PostJobsRequestOption {
	return func(r *PostJobsRequest) {
		r.GroupComment = &s
	}
}

func NewPostJobsRequest(jobs []*JobRequest, options ...PostJobsRequestOption) *PostJobsRequest {
	p := &PostJobsRequest{Jobs: jobs}
	for _, option := range options {
		option(p)
	}
	return p
}

type PostJobsResponse struct {
	OrderID     int               `json:"order_id"`
	Count       int               `json:"job_count"`
	CreditsUsed Float64           `json:"credits_used"`
	Currency    string            `json:"currency"`
	Jobs        []PostJobResponse `json:"jobs,omitempty"`
}

func (c *Client) PostJobs(req *PostJobsRequest) (*PostJobsResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	pjr := new(PostJobsResponse)
	err = c.post(jobsNamespace, bytes.NewReader(b), pjr)
	return pjr, err
}

type GetJobsRequest struct {
	Options url.Values
}

type GetJobsRequestOption func(*GetJobsRequest)

func WithStatus(s string) GetJobsRequestOption {
	return func(r *GetJobsRequest) {
		r.Options["status"] = []string{string(s)}
	}
}

func WithTimestampAfter(t Time) GetJobsRequestOption {
	return func(r *GetJobsRequest) {
		r.Options["timestamp_after"] = []string{strconv.Itoa(int(time.Time(t).Unix()))}
	}
}

func WithCount(i int) GetJobsRequestOption {
	return func(r *GetJobsRequest) {
		r.Options["count"] = []string{strconv.Itoa(i)}
	}
}

func NewGetJobsRequest(options ...GetJobsRequestOption) *GetJobsRequest {
	j := &GetJobsRequest{
		Options: url.Values{},
	}
	for _, option := range options {
		option(j)
	}
	return j
}

type GetJobsResponse struct {
	Jobs []GetJobResponse
}

func (g *GetJobsResponse) UnmarshalJSON(b []byte) error {
	j := new([]GetJobResponse)
	err := json.Unmarshal(b, j)
	if err != nil {
		return err
	}
	g.Jobs = *j
	return nil
}

func (c *Client) GetJobs(req *GetJobsRequest) (*GetJobsResponse, error) {
	pjr := new(GetJobsResponse)
	err := c.get(jobsNamespace, req.Options, pjr)
	return pjr, err
}

type GetJobsByIDRequest struct {
	IDs []int
}

func NewGetJobsByIDRequest(ids ...int) *GetJobsByIDRequest {
	return &GetJobsByIDRequest{IDs: ids}
}

type GetJobsByIDResponse struct {
	Jobs []GetJobResponse `json:"jobs,omitempty"`
}

func (c *Client) GetJobsByID(req *GetJobsByIDRequest) (*GetJobsByIDResponse, error) {
	strIDs := make([]string, len(req.IDs), len(req.IDs))
	for i := range req.IDs {
		strIDs[i] = strconv.Itoa(req.IDs[i])
	}
	pjr := new(GetJobsByIDResponse)
	err := c.get(fmt.Sprintf("%s/%s", jobsNamespace, strings.Join(strIDs, ",")), nil, pjr)
	return pjr, err
}

type ReviseJobsRequest struct {
	Action string              `json:"action"`
	Jobs   []*ReviseJobRequest `json:"job_ids"`
}

func NewReviseJobsRequest(jobs ...*ReviseJobRequest) *ReviseJobsRequest {
	return &ReviseJobsRequest{
		Action: "revise",
		Jobs:   jobs,
	}
}

func (c *Client) ReviseJobs(req *ReviseJobsRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = c.put(jobsNamespace, bytes.NewReader(b), nil)
	return err
}

type ArchiveJobsRequest struct {
	Action string               `json:"action"`
	Jobs   []*ArchiveJobRequest `json:"job_ids"`
}

func NewArchiveJobsRequest(jobs ...*ArchiveJobRequest) *ArchiveJobsRequest {
	return &ArchiveJobsRequest{
		Action: "archive",
		Jobs:   jobs,
	}
}

func (c *Client) ArchiveJobs(req *ArchiveJobsRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = c.put(jobsNamespace, bytes.NewReader(b), nil)
	return err
}

type ApproveJobsRequest struct {
	Action string               `json:"action"`
	Jobs   []*ApproveJobRequest `json:"job_ids"`
}

func NewApproveJobsRequest(jobs ...*ApproveJobRequest) *ApproveJobsRequest {
	return &ApproveJobsRequest{
		Action: "approve",
		Jobs:   jobs,
	}
}

func (c *Client) ApproveJobs(req *ApproveJobsRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = c.put(jobsNamespace, bytes.NewReader(b), nil)
	return err
}

type RejectJobsRequest struct {
	Action string              `json:"action"`
	Jobs   []*RejectJobRequest `json:"job_ids"`
}

func NewRejectJobsRequest(jobs ...*RejectJobRequest) *RejectJobsRequest {
	return &RejectJobsRequest{
		Action: "reject",
		Jobs:   jobs,
	}
}

type RejectJobsResponse struct {
	Jobs []RejectedJob `json:"jobs,omitempty"`
}

func (c *Client) RejectJobs(req *RejectJobsRequest) (*RejectJobsResponse, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	rjr := new(RejectJobsResponse)
	err = c.put(jobsNamespace, bytes.NewReader(b), rjr)
	return rjr, err
}
