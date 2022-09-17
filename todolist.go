package basecamp3

import "time"

type Todo struct {
	ID               int64     `json:"id"`
	Status           string    `json:"status"`
	VisibleToClients bool      `json:"visible_to_clients"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Title            string    `json:"title"`
	InheritsStatus   bool      `json:"inherits_status"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
	AppURL           string    `json:"app_url"`
	BookmarkURL      string    `json:"bookmark_url"`
	SubscriptionURL  string    `json:"subscription_url"`
	CommentsCount    int       `json:"comments_count"`
	CommentsURL      string    `json:"comments_url"`
	Position         int       `json:"position"`
	Parent           struct {
		ID     int64  `json:"id"`
		Title  string `json:"title"`
		Type   string `json:"type"`
		URL    string `json:"url"`
		AppURL string `json:"app_url"`
	} `json:"parent"`
	Bucket struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"bucket"`
	Creator struct {
		ID             int         `json:"id"`
		AttachableSgid string      `json:"attachable_sgid"`
		Name           string      `json:"name"`
		EmailAddress   string      `json:"email_address"`
		PersonableType string      `json:"personable_type"`
		Title          string      `json:"title"`
		Bio            string      `json:"bio"`
		Location       interface{} `json:"location"`
		CreatedAt      time.Time   `json:"created_at"`
		UpdatedAt      time.Time   `json:"updated_at"`
		Admin          bool        `json:"admin"`
		Owner          bool        `json:"owner"`
		Client         bool        `json:"clietime"`
		Employee       bool        `json:"employee"`
		TimeZone       string      `json:"time_zone"`
		AvatarURL      string      `json:"avatar_url"`
		Company        struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"company"`
	} `json:"creator"`
	Description    string `json:"description"`
	Completed      bool   `json:"completed"`
	CompletedRatio string `json:"completed_ratio"`
	Name           string `json:"name"`
	TodosURL       string `json:"todos_url"`
	GroupsURL      string `json:"groups_url"`
	AppTodosURL    string `json:"app_todos_url"`
}
