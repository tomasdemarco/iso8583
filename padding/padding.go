package padding

type Padding struct {
	Type     Type     `json:"type"`
	Position Position `json:"position"`
	Char     string   `json:"char"`
}
