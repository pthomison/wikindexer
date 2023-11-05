package types

import "gitlab.com/tozd/go/mediawiki"

type Article struct {
	Id       int64
	Name     string
	Abstract string
	URL      string
	WikiText string
}

func ParseArticle(raw mediawiki.Article) *Article {
	a := &Article{
		Id:       raw.Identifier,
		Name:     raw.Name,
		Abstract: raw.Abstract,
		URL:      raw.URL,
		WikiText: raw.ArticleBody.WikiText,
	}

	return a
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

func (a *Article) BatchInsert() string {
	s := `INSERT OR IGNORE INTO article (id, name, url, abstract)
		VALUES (:id, :name, :url, :abstract)`

	return s
}
