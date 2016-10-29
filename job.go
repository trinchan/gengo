package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/trinchan/gengo/lang"
)

const (
	jobNamespace = "/translate/job"
)

type PreferredTranslator struct {
	ID        int  `json:"id"`
	LastLogin Time `json:"last_login"`
}

const (
	JobTypeText = "text"
	JobTypeFile = "file"
)

type Attachment struct {
	URL      string `json:"url"`
	Name     string `json:"filename"`
	MimeType string `json:"mime_type"`
}

type Tier string

const (
	TierStandard = Tier("standard")
	TierPro      = Tier("pro")
	TierUltra    = Tier("ultra")
)

type JobRequest struct {
	Type string `json:"type"`
	lang.Pair
	Tier         Tier         `json:"tier"`
	Slug         string       `json:"slug"`
	Attachments  []Attachment `json:"attachments"`
	BodySrc      *string      `json:"body_src,omitempty"`
	Force        Bool         `json:"force,omitempty"`
	Comment      string       `json:"comment,omitempty"`
	UsePreferred Bool         `json:"use_preferred,omitempty"`
	CallbackURL  string       `json:"callback_url,omitempty"`
	AutoApprove  Bool         `json:"auto_approve,omitempty"`
	CustomData   string       `json:"custom_data,omitempty"`
	Purpose      string       `json:"purpose,omitempty"`
	AsGroup      Bool         `json:"as_group,omitempty"`
	GlossaryID   *int         `json:"glossary_id,omitempty"`
	MaxChars     *int         `json:"max_chars,omitempty"`
	Position     *int         `json:"position,omitempty"`
	Identifier   *string      `json:"identifier,omitempty"`
}

type FileJobRequest struct {
	*JobRequest
	FilePath string `json:"-"`
	FileKey  string `json:"file_key"`
}

func NewFileJobRequest(filename string, lp lang.Pair, tier Tier, options ...JobOption) *FileJobRequest {
	jr := &JobRequest{
		Type: "file",
		Pair: lp,
		Tier: tier,
	}
	for _, option := range options {
		option(jr)
	}
	fjr := &FileJobRequest{
		JobRequest: jr,
		FilePath:   filename,
	}
	return fjr
}

func NewJobRequest(text string, lp lang.Pair, tier Tier, options ...JobOption) *JobRequest {
	jr := &JobRequest{
		Type:    "text",
		Pair:    lp,
		Tier:    tier,
		BodySrc: &text,
	}
	for _, option := range options {
		option(jr)
	}
	return jr
}

type JobOption func(*JobRequest)

func WithSlug(s string) JobOption {
	return func(jr *JobRequest) {
		jr.Slug = s
	}
}

func WithForce(b Bool) JobOption {
	return func(jr *JobRequest) {
		jr.Force = b
	}
}

func WithComment(c string) JobOption {
	return func(jr *JobRequest) {
		jr.Comment = c
	}
}

func WithPreferred(b Bool) JobOption {
	return func(jr *JobRequest) {
		jr.UsePreferred = b
	}
}

func WithAttachments(a []Attachment) JobOption {
	return func(jr *JobRequest) {
		jr.Attachments = a
	}
}

func WithCallbackURL(s string) JobOption {
	return func(jr *JobRequest) {
		jr.CallbackURL = s
	}
}

func WithAutoApprove(b Bool) JobOption {
	return func(jr *JobRequest) {
		jr.AutoApprove = b
	}
}

func WithCustomData(s string) JobOption {
	return func(jr *JobRequest) {
		jr.CustomData = s
	}
}

func WithPurpose(s string) JobOption {
	return func(jr *JobRequest) {
		jr.Purpose = s
	}
}

func WithAsGroup(b Bool) JobOption {
	return func(jr *JobRequest) {
		jr.AsGroup = b
	}
}

func WithGlossaryID(i int) JobOption {
	return func(jr *JobRequest) {
		jr.GlossaryID = &i
	}
}

func WithMaxChars(i int) JobOption {
	return func(jr *JobRequest) {
		jr.MaxChars = &i
	}
}

func WithPosition(i int) JobOption {
	return func(jr *JobRequest) {
		jr.Position = &i
	}
}

type PostJobResponse struct {
	ID      Int    `json:"job_id"`
	OrderID Int    `json:"order_id"`
	BodySrc string `json:"body_src"`
	BodyTgt string `json:"body_tgt"`
	lang.Pair
	Tier
	UnitCount          Int     `json:"unit_count"`
	Credits            Float64 `json:"credits"`
	Status             string  `json:"status"`
	CaptchaURL         string  `json:"captcha_url"`
	ETA                int     `json:"eta"`
	CallbackURL        string  `json:"callback_url"`
	AutoApprove        Bool    `json:"auto_approve"`
	Ctime              Time    `json:"ctime"`
	CustomData         string  `json:"custom_data"`
	MachineTranslation Bool    `json:"mt"`
	FileSourceURL      string  `json:"file_url_src"`
	FileTargetURL      string  `json:"file_url_tgt"`
}

func (jr *PostJobResponse) UnmarshalJSON(d []byte) error {
	// hack off the double array container cuz gengo api is wonky
	type c PostJobResponse
	var x = new(c)
	err := json.Unmarshal(d[1:len(d)-1], x)
	if err != nil {
		return err
	}
	*jr = PostJobResponse(*x)
	return nil
}

type GetJobResponse PostJobResponse

type ReviseJobRequest struct {
	ID      int     `json:"job_id"`
	Comment *string `json:"comment,omitempty"`
}

func NewReviseJobRequest(id int, options ...ReviseJobOption) *ReviseJobRequest {
	rjr := &ReviseJobRequest{
		ID: id,
	}
	for _, option := range options {
		option(rjr)
	}
	return rjr
}

type ReviseJobOption func(*ReviseJobRequest)

func WithRevisionComment(s string) ReviseJobOption {
	return func(r *ReviseJobRequest) {
		r.Comment = &s
	}
}

type ApproveJobRequest struct {
	ID                   int     `json:"job_id"`
	Rating               *int    `json:"rating,omitempty"`
	CommentForTranslator *string `json:"for_translator,omitempty"`
	CommentForGengo      *string `json:"for_mygengo,omitempty"`
	Public               *Bool   `json:"public,omitempty"`
}

func NewApproveJobRequest(id int, options ...ApproveJobOption) *ApproveJobRequest {
	rjr := &ApproveJobRequest{
		ID: id,
	}
	for _, option := range options {
		option(rjr)
	}
	return rjr
}

type ApproveJobOption func(*ApproveJobRequest)

func WithRating(i int) ApproveJobOption {
	return func(r *ApproveJobRequest) {
		r.Rating = &i
	}
}

func WithTranslatorComment(s string) ApproveJobOption {
	return func(r *ApproveJobRequest) {
		r.CommentForTranslator = &s
	}
}

func WithGengoComment(s string) ApproveJobOption {
	return func(r *ApproveJobRequest) {
		r.CommentForGengo = &s
	}
}

func WithPublicComment(b Bool) ApproveJobOption {
	return func(r *ApproveJobRequest) {
		r.Public = &b
	}
}

type RejectedJob struct {
	ID           int    `json:"job_id"`
	CustomerID   int    `json:"customer_id"`
	TranslatorID int    `json:"worker_id"`
	Comment      string `json:"rejection_comments"`
	Reason       string `json:"rejection_reason"`
}

type RejectJobRequest struct {
	ID         int     `json:"job_id"`
	Reason     string  `json:"reason"`
	Comment    string  `json:"comment"`
	Captcha    string  `json:"captcha"`
	CaptchaURL *string `json:"captcha_url,omitempty"`
	FollowUp   *string `json:"follow_up,omitempty"`
}

const (
	RejectionReasonQuality    = "quality"
	RejectionReasonIncomplete = "incomplete"
	RejectionReasonOther      = "other"
)

func NewRejectJobRequest(id int, reason, comment, captcha string, options ...RejectJobOption) *RejectJobRequest {
	rjr := &RejectJobRequest{
		ID:      id,
		Reason:  reason,
		Comment: comment,
		Captcha: captcha,
	}
	for _, option := range options {
		option(rjr)
	}
	return rjr
}

type RejectJobOption func(*RejectJobRequest)

func WithFollowUp(s string) RejectJobOption {
	return func(r *RejectJobRequest) {
		r.FollowUp = &s
	}
}

func WithCaptchaURL(url string) RejectJobOption {
	return func(r *RejectJobRequest) {
		r.CaptchaURL = &url
	}
}

type ArchiveJobRequest int

func NewArchiveJobRequest(id int) *ArchiveJobRequest {
	ajr := ArchiveJobRequest(id)
	return &ajr
}

type GetJobRequest struct {
	ID int
}

func NewGetJobRequest(id int) *GetJobRequest {
	return &GetJobRequest{ID: id}
}

type GetJobByIDResponse struct {
	Job GetJobResponse `json:"job,omitempty"`
}

func (c *Client) GetJob(req *GetJobRequest) (*GetJobByIDResponse, error) {
	pjr := new(GetJobByIDResponse)
	err := c.get(jobNamespace+fmt.Sprintf("/%d", req.ID), nil, pjr)
	return pjr, err
}

type CancelJobRequest struct {
	ID int
}

func NewCancelJobRequest(id int) *CancelJobRequest {
	return &CancelJobRequest{ID: id}
}

func (c *Client) CancelJob(req *CancelJobRequest) error {
	err := c.delete(jobNamespace+fmt.Sprintf("/%d", req.ID), nil, nil)
	return err
}

type JobRevisionsRequest struct {
	ID int
}

func NewJobRevisionsRequest(id int) *JobRevisionsRequest {
	return &JobRevisionsRequest{ID: id}
}

type JobRevisionsResponse struct {
	JobID     int              `json:"job_id"`
	Revisions []RevisionWithID `json:"revisions"`
}

type RevisionWithID struct {
	ID    int  `json:"rev_id"`
	Ctime Time `json:"ctime"`
}

func (c *Client) JobRevisions(req *JobRevisionsRequest) (*JobRevisionsResponse, error) {
	gjr := new(JobRevisionsResponse)
	err := c.get(jobNamespace+fmt.Sprintf("/%d/revisions", req.ID), nil, gjr)
	return gjr, err
}

type JobRevisionRequest struct {
	ID         int
	RevisionID int
}

func NewJobRevisionRequest(id, revisionID int) *JobRevisionRequest {
	return &JobRevisionRequest{ID: id, RevisionID: revisionID}
}

type JobRevisionResponse struct {
	Revision RevisionWithBody `json:"revision"`
}

type RevisionWithBody struct {
	Body  string `json:"body_tgt"`
	Ctime Time   `json:"ctime"`
}

func (c *Client) JobRevision(req *JobRevisionRequest) (*JobRevisionResponse, error) {
	gjr := new(JobRevisionResponse)
	err := c.get(jobNamespace+fmt.Sprintf("/%d/revisions/%d", req.ID, req.RevisionID), nil, gjr)
	return gjr, err
}

type JobFeedbackRequest struct {
	ID int
}

func NewJobFeedbackRequest(id int) *JobFeedbackRequest {
	return &JobFeedbackRequest{ID: id}
}

type JobFeedbackResponse struct {
	Feedback Feedback `json:"feedback"`
}

type Feedback struct {
	Comment string `json:"for_translator"`
	Rating  int    `json:"rating"`
}

func (c *Client) JobFeedback(req *JobFeedbackRequest) (*JobFeedbackResponse, error) {
	gjr := new(JobFeedbackResponse)
	err := c.get(jobNamespace+fmt.Sprintf("/%d/feedback", req.ID), nil, gjr)
	return gjr, err
}

type JobCommentsRequest struct {
	ID int
}

func NewJobCommentsRequest(id int) *JobCommentsRequest {
	return &JobCommentsRequest{ID: id}
}

type JobCommentsResponse CommentsResponse

func (c *Client) JobComments(req *JobCommentsRequest) (*JobCommentsResponse, error) {
	gjr := new(JobCommentsResponse)
	err := c.get(jobNamespace+fmt.Sprintf("/%d/comments", req.ID), nil, gjr)
	return gjr, err
}

type AddJobCommentRequest struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewAddJobCommentRequest(id int, comment string) *AddJobCommentRequest {
	r := &AddJobCommentRequest{
		ID:   id,
		Body: comment,
	}
	return r
}

func (c *Client) AddJobComment(req *AddJobCommentRequest) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = c.post(jobNamespace+fmt.Sprintf("/%d/comment", req.ID), bytes.NewReader(b), nil)
	return err
}
