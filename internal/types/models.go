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

type LatestTweetStats struct {
	Likes    int  `json:"likes"`
	Retweets int  `json:"retweets"`
	Replies  int  `json:"replies"`
	Ratioed  bool `json:"ratioed"`
}

type TwitterProfile struct {
	Handle           string           `json:"handle"`
	DisplayName      string           `json:"display_name"`
	VerifiedReason   string           `json:"verified_reason"`
	Bio              string           `json:"bio"`
	Location         string           `json:"location"`
	Followers        int              `json:"followers"`
	Following        int              `json:"following"`
	Joined           string           `json:"joined"`
	ProfileImage     string           `json:"profile_image"`
	BannerCaption    string           `json:"banner_caption"`
	PinnedTweet      string           `json:"pinned_tweet"`
	PinnedSkills     []string         `json:"pinned_skills"`
	Tweets           []string         `json:"tweets"`
	Endorsements     []string         `json:"endorsements"`
	SpacesHosted     []string         `json:"spaces_hosted"`
	Replies          []string         `json:"replies"`
	LatestTweetStats LatestTweetStats `json:"latest_tweet_stats"`
	Score            int              `json:"score"`
}
