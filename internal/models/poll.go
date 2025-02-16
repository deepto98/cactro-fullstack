package models

// Poll represents a poll with a question and a list of options.
type Poll struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []Option `json:"options"`
}

// Option represents a poll option with a vote count.
type Option struct {
	ID        int    `json:"id"`
	PollID    int    `json:"poll_id"`
	Text      string `json:"option_text"`
	VoteCount int    `json:"vote_count"`
}
