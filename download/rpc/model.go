package rpc

// Parameter contains all request actions
type Parameter struct {
	Website string   `json:"website"`
	Comic   string   `json:"comic"`
	Chapter string   `json:"chapter"`
	Path    string   `json:"path"`
	Latest  string   `json:"latest"`
	URLs    []string `json:"urls"`
}

// Comic contains comic informations
type Comic struct {
	URL      string            `json:"url"`
	Title    string            `json:"title"`
	Summary  string            `json:"summary"`
	Cover    string            `json:"cover"`
	Chapters map[string]string `json:"chapters"`
	Indexes  []string          `json:"indexes"`
	Latest   string            `json:"latest"`
	Source   string            `json:"source"`
}

// Chapter contains chapter informations
type Chapter struct {
	URL      string            `json:"url"`
	Title    string            `json:"title"`
	Ctitle   string            `json:"ctitle"`
	Pictures map[string]string `json:"pictures"`
}

// Update contains updated information
type Update struct {
	Chapters map[string]string `json:"chapters"`
	Indexes  []string          `json:"indexes"`
	Latest   string            `json:"latest"`
}
