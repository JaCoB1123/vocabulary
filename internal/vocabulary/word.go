package vocabulary

type WordPair struct {
	Name        string
	Translation string
	Attributes  map[string]string `json:",omitempty"`
	Tags        []string
}
