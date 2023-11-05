package types

type Article struct {
	Id       int64
	Name     string
	Abstract string
	URL      string
	WikiText string
}

func (a *Article) Schema() string {
	s := `
CREATE TABLE IF NOT EXISTS article (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	abstract TEXT,
	wiki_text TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_article_id 
ON article (id);
	`
	return s
}
