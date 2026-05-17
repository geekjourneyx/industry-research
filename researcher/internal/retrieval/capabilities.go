package retrieval

type ProviderCapabilities struct {
	Provider               string       `json:"provider"`
	ProviderType           ProviderType `json:"provider_type"`
	Modes                  []Mode       `json:"modes"`
	SupportsFreshness      bool         `json:"supports_freshness"`
	SupportsIncludeDomains bool         `json:"supports_include_domains"`
	SupportsExcludeDomains bool         `json:"supports_exclude_domains"`
	SupportsSummary        bool         `json:"supports_summary"`
	SupportsLocation       bool         `json:"supports_location"`
	SupportsSources        bool         `json:"supports_sources"`
	SupportsImages         bool         `json:"supports_images"`
	SupportsModelChoice    bool         `json:"supports_model_choice"`
	ResultKinds            []string     `json:"result_kinds"`
}

func BuiltInCapabilities() []ProviderCapabilities {
	return []ProviderCapabilities{
		{
			Provider:               "bocha",
			ProviderType:           ProviderTypeDirectSearch,
			Modes:                  []Mode{ModeSearch, ModeRetrieve},
			SupportsFreshness:      true,
			SupportsIncludeDomains: true,
			SupportsExcludeDomains: true,
			SupportsSummary:        true,
			SupportsLocation:       false,
			SupportsSources:        false,
			SupportsImages:         true,
			SupportsModelChoice:    false,
			ResultKinds:            []string{"web_page", "image"},
		},
		{
			Provider:               "volcengine",
			ProviderType:           ProviderTypeModelAnswerSearch,
			Modes:                  []Mode{ModeAnswer, ModeRetrieve},
			SupportsFreshness:      false,
			SupportsIncludeDomains: false,
			SupportsExcludeDomains: false,
			SupportsSummary:        false,
			SupportsLocation:       true,
			SupportsSources:        true,
			SupportsImages:         false,
			SupportsModelChoice:    true,
			ResultKinds:            []string{"annotation_url", "answer_text", "retrieval_call"},
		},
	}
}
