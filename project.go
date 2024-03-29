package basecamp3

import "time"

type Project struct {
	ID             int       `json:"id"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Purpose        string    `json:"purpose"`
	ClientsEnabled bool      `json:"clients_enabled"`
	BookmarkURL    string    `json:"bookmark_url"`
	URL            string    `json:"url"`
	AppURL         string    `json:"app_url"`
	Dock           []*Dock   `json:"dock"`
}

type Dock struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Position int    `json:"position"`
	URL      string `json:"url"`
	AppURL   string `json:"app_url"`
}

func (p *Project) GetDock(dockname string) *Dock {
	for _, dock := range p.Dock {
		if dock.Name == dockname {
			return dock
		}
	}
	return nil
}

func (p *Project) GetTodoSet() int64 {
	d := p.GetDock("todoset")
	if d == nil {
		return -1
	}
	return d.ID
}
