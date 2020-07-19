package vocabulary

type WordPair struct {
	Name        string
	Translation string
	Attributes  map[string]string `json:",omitempty"`
	Tags        []string
}

func (word WordPair) IsFilteredBy(tags []string) bool {
	return !containsAll(tags, word.Tags)
}
