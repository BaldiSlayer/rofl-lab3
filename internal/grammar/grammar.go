package grammar

import "github.com/BaldiSlayer/rofl-lab3/internal/models"

type Grammar struct {
}

func New(rules []models.Rule) *Grammar {
	return &Grammar{}
}
