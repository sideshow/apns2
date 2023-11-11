package payload

type aps struct {
	Alert             interface{}            `json:"alert,omitempty"`
	Badge             interface{}            `json:"badge,omitempty"`
	Sound             interface{}            `json:"sound,omitempty"`
	ThreadID          string                 `json:"thread-id,omitempty"`
	Category          string                 `json:"category,omitempty"`
	ContentAvailable  int                    `json:"content-available,omitempty"`
	MutableContent    int                    `json:"mutable-content,omitempty"`
	InterruptionLevel EInterruptionLevel     `json:"interruption-level,omitempty"`
	RelevanceScore    interface{}            `json:"relevance-score,omitempty"`
	StaleDate         int64                  `json:"stale-date,omitempty"`
	ContentState      map[string]interface{} `json:"content-state,omitempty"`
	Timestamp         int64                  `json:"timestamp,omitempty"`
	Event             ELiveActivityEvent     `json:"event,omitempty"`
	DismissalDate     int64                  `json:"dismissal-date,omitempty"`
	URLArgs           []string               `json:"url-args,omitempty"`
}
