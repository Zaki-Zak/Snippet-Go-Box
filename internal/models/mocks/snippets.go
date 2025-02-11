package mocks

import (
	"time"

	"github.com/Zaki-Zak/Snippet-Go-Box/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "little pond",
	Content: "little pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (models.Snippet, error) {
	if id == 1 {
		return mockSnippet, nil
	} else {
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
