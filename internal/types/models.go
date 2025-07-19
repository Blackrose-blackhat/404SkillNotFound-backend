package types

type JudgeOutput struct {
	Score          int            `json:"score"`
	Grade          Grade          `json:"grade"` // now a nested struct
	Summary        string         `json:"summary"`
	Feedback       []Feedback     `json:"feedback"`
	RoastMode      bool           `json:"roast_mode"`
	Recommendation Recommendation `json:"recommendation"`
	Reaction       string         `json:"reaction"`
}

type Grade struct {
	Letter      string `json:"letter"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type Feedback struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type Recommendation struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
