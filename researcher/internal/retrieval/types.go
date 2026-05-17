package retrieval

import "time"

type ProviderType string

const (
	ProviderTypeDirectSearch        ProviderType = "direct_search"
	ProviderTypeModelAnswerSearch   ProviderType = "model_answer_search"
	ProviderTypeAgentNativeSearch   ProviderType = "agent_native_search"
	ProviderTypeBrowserVerification ProviderType = "browser_verification"
	ProviderTypeKnowledgeRetrieval  ProviderType = "knowledge_retrieval"
)

type Mode string

const (
	ModeSearch   Mode = "search"
	ModeAnswer   Mode = "answer"
	ModeRetrieve Mode = "retrieve"
)

type RetrievalRequest struct {
	Provider     string         `json:"provider"`
	ProviderType ProviderType   `json:"provider_type"`
	Mode         Mode           `json:"mode"`
	Query        string         `json:"query"`
	Parameters   map[string]any `json:"parameters"`
	Headers      map[string]any `json:"-"`
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
	DisplayURL           string         `json:"display_url"`
	SiteName             string         `json:"site_name"`
	SiteIcon             string         `json:"site_icon"`
	Snippet              string         `json:"snippet"`
	Summary              string         `json:"summary"`
	PublishedAt          time.Time      `json:"published_at"`
	LastCrawledAt        time.Time      `json:"last_crawled_at"`
	Language             string         `json:"language"`
	ContentType          string         `json:"content_type"`
	SourceConfidenceHint string         `json:"source_confidence_hint"`
	ProviderMetadata     map[string]any `json:"provider_metadata"`
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
	ProviderStatus int    `json:"provider_status"`
	ProviderCode   string `json:"provider_code"`
	ProviderLogID  string `json:"provider_log_id"`
	Retryable      bool   `json:"retryable"`
	AgentAction    string `json:"agent_action"`
	RawErrorPath   string `json:"raw_error_path"`
}
