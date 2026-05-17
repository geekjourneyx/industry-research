package retrieval

import "time"

type ProviderType string

const (
	ProviderTypeDirectSearch       ProviderType = "direct_search"
	ProviderTypeModelAnswerSearch  ProviderType = "model_answer_search"
	ProviderTypeAgentNativeSearch  ProviderType = "agent_native_search"
	ProviderTypeBrowserVerify      ProviderType = "browser_verification"
	ProviderTypeKnowledgeRetrieval ProviderType = "knowledge_retrieval"
	ProviderTypeMulti              ProviderType = "multi"
)

type Mode string

const (
	ModeSearch   Mode = "search"
	ModeAnswer   Mode = "answer"
	ModeRetrieve Mode = "retrieve"
)

type RetrievalRequest struct {
	Provider     string            `json:"provider"`
	ProviderType ProviderType      `json:"provider_type,omitempty"`
	Mode         Mode              `json:"mode"`
	Query        string            `json:"query"`
	Parameters   map[string]any    `json:"parameters,omitempty"`
	Headers      map[string]string `json:"-"`
}

type RetrievalResponse struct {
	Provider       string          `json:"provider"`
	ProviderType   ProviderType    `json:"provider_type"`
	Mode           Mode            `json:"mode"`
	Query          string          `json:"query"`
	RetrievedAt    time.Time       `json:"retrieved_at"`
	Request        map[string]any  `json:"request"`
	RetrievalCalls []RetrievalCall `json:"retrieval_calls"`
	Items          []Item          `json:"items"`
	Answer         Answer          `json:"answer"`
	Usage          map[string]any  `json:"usage"`
	Errors         []Error         `json:"errors"`
}

type MultiResponse struct {
	Provider        string              `json:"provider"`
	ProviderType    ProviderType        `json:"provider_type"`
	Mode            Mode                `json:"mode"`
	Query           string              `json:"query"`
	RetrievedAt     time.Time           `json:"retrieved_at"`
	ProviderResults []RetrievalResponse `json:"provider_results"`
	Errors          []Error             `json:"errors"`
}

type RetrievalCall struct {
	CallID           string         `json:"call_id"`
	Query            string         `json:"query"`
	Status           string         `json:"status"`
	ProviderAction   string         `json:"provider_action"`
	ProviderMetadata map[string]any `json:"provider_metadata"`
}

type Item struct {
	Rank                 int            `json:"rank"`
	Title                string         `json:"title"`
	URL                  string         `json:"url"`
	DisplayURL           string         `json:"display_url,omitempty"`
	SiteName             string         `json:"site_name,omitempty"`
	SiteIcon             string         `json:"site_icon,omitempty"`
	Snippet              string         `json:"snippet,omitempty"`
	Summary              string         `json:"summary,omitempty"`
	PublishedAt          string         `json:"published_at,omitempty"`
	LastCrawledAt        string         `json:"last_crawled_at,omitempty"`
	Language             string         `json:"language,omitempty"`
	ContentType          string         `json:"content_type"`
	SourceConfidenceHint string         `json:"source_confidence_hint"`
	ProviderMetadata     map[string]any `json:"provider_metadata,omitempty"`
}

type Answer struct {
	Text      string     `json:"text"`
	Citations []Citation `json:"citations"`
}

type Citation struct {
	Index  int    `json:"index"`
	URL    string `json:"url"`
	Title  string `json:"title"`
	Source string `json:"source"`
}

type Error struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	ProviderStatus int    `json:"provider_status,omitempty"`
	ProviderCode   string `json:"provider_code,omitempty"`
	ProviderLogID  string `json:"provider_log_id,omitempty"`
	Retryable      bool   `json:"retryable"`
	AgentAction    string `json:"agent_action"`
	RawErrorPath   string `json:"raw_error_path,omitempty"`
}
