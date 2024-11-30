package todo_test

import (
	"context"
	"my-first-api/internal/db"
	"my-first-api/internal/todo"
	"reflect"
	"testing"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {

	m.items = append(m.items, item)
	return nil
}

func (m MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil
}

func TestService_Search(t *testing.T) {
	tests := []struct {
		name           string
		toDosToAdd     []string
		query          string
		expectedResult []string
	}{
		{
			name:           "given a todo of gas and a search of ga, I should get a gas station back",
			toDosToAdd:     []string{"gas"},
			query:          "ga",
			expectedResult: []string{"gas"},
		},
		{
			name:           "still returns gas, even if case doesn't match",
			toDosToAdd:     []string{"Gas"},
			query:          "ga",
			expectedResult: []string{"Gas"},
		},
		{
			name:           "spaces",
			toDosToAdd:     []string{"go to gas"},
			query:          "go",
			expectedResult: []string{"go to gas"},
		},
		{
			name:           "space at start of word",
			toDosToAdd:     []string{" Space at beginning"},
			query:          "space",
			expectedResult: []string{" Space at beginning"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}
			svc := todo.NewService(m)
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}

			}
			got, err := svc.Search(tt.query)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
