package ledger

import "testing"

func TestSum(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]int
		want    int
	}{
		{
			name:    "empty store",
			entries: nil,
			want:    0,
		},
		{
			name:    "single value",
			entries: map[string]int{"a": 42},
			want:    42,
		},
		{
			name:    "multiple values",
			entries: map[string]int{"a": 1, "b": 2, "c": 3},
			want:    6,
		},
		{
			name:    "negative values",
			entries: map[string]int{"a": 10, "b": -4, "c": -1},
			want:    5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStore()
			for k, v := range tt.entries {
				s.Set(k, v)
			}
			if got := s.Sum(); got != tt.want {
				t.Errorf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}
