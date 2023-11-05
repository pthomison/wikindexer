package types

type Article struct {
	Name     string
	Abstract string
	URL      string
	WikiText string
}

func (a *Article) Schema() string {
	s := `
CREATE TABLE IF NOT EXISTS article (
	name text,
	abstract text,
	url text,
	wiki_text text
);
	`
	return s
}
