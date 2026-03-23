package domain_test

import (
	"testing"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/authors/domain"
	"github.com/stretchr/testify/assert"
)

func TestBlogCountChangedEvent_EventName(t *testing.T) {
	tests := []struct {
		name        string
		isIncrement bool
		expected    string
	}{
		{"increment", true, "authorIdentity.blogCountIncreased"},
		{"decrement", false, "authorIdentity.blogCountDecreased"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := domain.BlogCountChangedEvent{
				IsIncrement: tt.isIncrement,
			}

			assert.Equal(t, tt.expected, event.EventName())
		})
	}
}
