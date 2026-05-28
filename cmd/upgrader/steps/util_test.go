package steps

import (
	"strings"
	"testing"
)

func TestHelmTrackedAddonsSettled(t *testing.T) {
	tests := []struct {
		name        string
		states      []helmTrackedAddonRuntime
		wantSettled bool
		wantErr     string
	}{
		{
			name: "enabled addon is settled",
			states: []helmTrackedAddonRuntime{
				{
					Name:               "mysql",
					Phase:              "Enabled",
					Generation:         2,
					ObservedGeneration: 2,
					JobState:           "Succeeded",
				},
			},
			wantSettled: true,
		},
		{
			name: "failed addon returns error",
			states: []helmTrackedAddonRuntime{
				{
					Name:               "redis",
					Phase:              "Failed",
					Generation:         2,
					ObservedGeneration: 2,
					JobState:           "Failed",
				},
			},
			wantErr: "addon 进入 Failed 状态: redis",
		},
		{
			name: "addon with active job is not settled",
			states: []helmTrackedAddonRuntime{
				{
					Name:               "redis",
					Phase:              "Enabled",
					Generation:         2,
					ObservedGeneration: 2,
					JobState:           "Active",
				},
			},
		},
		{
			name: "addon waiting for observed generation is not settled",
			states: []helmTrackedAddonRuntime{
				{
					Name:               "redis",
					Phase:              "Enabled",
					Generation:         2,
					ObservedGeneration: 1,
					JobState:           "Succeeded",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSettled, err := helmTrackedAddonsSettled(tt.states)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotSettled != tt.wantSettled {
				t.Fatalf("settled = %v, want %v", gotSettled, tt.wantSettled)
			}
		})
	}
}
