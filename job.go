package gengo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/trinchan/gengo/language"
)

const (
	jobNamespace = "/translate/job"
)

type PreferredTranslator struct {
	ID        int   `json:"id"`
	LastLogin int64 `json:"last_login"`
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
	language.Pair
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

func NewFileJobRequest(filename string, lp language.Pair, tier Tier, options ...JobOption) *FileJobRequest {
	jr := &JobRequest{
		Type: "file",
		Pair: lp,
		Tier: tier,
	}
	for _, option := range options {
		option.Apply(jr)
	}
	fjr := &FileJobRequest{
		JobRequest: jr,
		FilePath:   filename,
	}
	return fjr
}

func NewJobRequest(text string, lp language.Pair, tier Tier, options ...JobOption) *JobRequest {
	jr := &JobRequest{
		Type:    "text",
		Pair:    lp,
		Tier:    tier,
		BodySrc: &text,
	}
	for _, option := range options {
		option.Apply(jr)
	}
	return jr
}

type JobOption interface {
	Apply(*JobRequest)
}

type SlugOption string

func (o SlugOption) Apply(jr *JobRequest) {
	jr.Slug = string(o)
}

type ForceOption Bool

func (o ForceOption) Apply(jr *JobRequest) {
	jr.Force = Bool(o)
}

type CommentOption struct {
	Comment string
}

func (o CommentOption) Apply(jr *JobRequest) {
	jr.Comment = o.Comment
}

type PreferredOption Bool

func (o PreferredOption) Apply(jr *JobRequest) {
	jr.UsePreferred = Bool(o)
}

type AttachmentOption []Attachment

func (o AttachmentOption) Apply(jr *JobRequest) {
	jr.Attachments = []Attachment(o)
}

type CallbackURLOption string

func (o CallbackURLOption) Apply(jr *JobRequest) {
	jr.CallbackURL = string(o)
}

type AutoApproveOption Bool

func (o AutoApproveOption) Apply(jr *JobRequest) {
	jr.AutoApprove = Bool(o)
}

type CustomDataOption string

func (o CustomDataOption) Apply(jr *JobRequest) {
	jr.CustomData = string(o)
}

type PurposeOption string

func (o PurposeOption) Apply(jr *JobRequest) {
	jr.Purpose = string(o)
}

type GroupOption Bool

func (o GroupOption) Apply(jr *JobRequest) {
	jr.AsGroup = Bool(o)
}

type GlossaryOption int

func (o GlossaryOption) Apply(jr *JobRequest) {
	id := int(o)
	jr.GlossaryID = &id
}

type MaxCharsOption int

func (o MaxCharsOption) Apply(jr *JobRequest) {
	max := int(o)
	jr.MaxChars = &max
}

type PositionOption int

func (o PositionOption) Apply(jr *JobRequest) {
	p := int(o)
	jr.Position = &p
}

type PostJobResponse struct {
	ID      Int    `json:"job_id"`
	OrderID Int    `json:"order_id"`
	BodySrc string `json:"body_src"`
	BodyTgt string `json:"body_tgt"`
	language.Pair
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
		option.Apply(rjr)
	}
	return rjr
}

type ReviseJobOption interface {
	Apply(*ReviseJobRequest)
}

type ReviseJobCommentOption string

func (o ReviseJobCommentOption) Apply(r *ReviseJobRequest) {
	c := string(o)
	r.Comment = &c
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
		option.Apply(rjr)
	}
	return rjr
}

type ApproveJobOption interface {
	Apply(*ApproveJobRequest)
}

type ApproveJobRatingOption int

func (o ApproveJobRatingOption) Apply(r *ApproveJobRequest) {
	rating := int(o)
	r.Rating = &rating
}

type ApproveJobTranslatorCommentOption string

func (o ApproveJobTranslatorCommentOption) Apply(r *ApproveJobRequest) {
	c := string(o)
	r.CommentForTranslator = &c
}

type ApproveJobGengoCommentOption string

func (o ApproveJobGengoCommentOption) Apply(r *ApproveJobRequest) {
	c := string(o)
	r.CommentForGengo = &c
}

type ApproveJobPublicCommentOption Bool

func (o ApproveJobPublicCommentOption) Apply(r *ApproveJobRequest) {
	p := Bool(o)
	r.Public = &p
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
		option.Apply(rjr)
	}
	return rjr
}

type RejectJobOption interface {
	Apply(*RejectJobRequest)
}

type RejectJobFollowUpOption string

func (o RejectJobFollowUpOption) Apply(r *RejectJobRequest) {
	f := string(o)
	r.FollowUp = &f
}

type RejectJobCaptchaURLOption string

func (o RejectJobCaptchaURLOption) Apply(r *RejectJobRequest) {
	c := string(o)
	r.CaptchaURL = &c
}

type ArchiveJobRequest int

func NewArchiveJobRequest(id int, options ...ArchiveJobOption) *ArchiveJobRequest {
	ajr := ArchiveJobRequest(id)
	for _, option := range options {
		option.Apply(&ajr)
	}
	return &ajr
}

type ArchiveJobOption interface {
	Apply(*ArchiveJobRequest)
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

type AddJobCommentRequestOption interface {
	Apply(*AddJobCommentRequest)
}

func NewAddJobCommentRequest(id int, comment string, options ...AddJobCommentRequestOption) *AddJobCommentRequest {
	r := &AddJobCommentRequest{
		ID:   id,
		Body: comment,
	}
	for _, option := range options {
		option.Apply(r)
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
