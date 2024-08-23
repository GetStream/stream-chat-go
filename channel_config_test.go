package stream_chat

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input DurationString
		want  string
	}{
		{
			name:  "Zero",
			input: DurationString(0),
			want:  `null`,
		},
		{
			name:  "Hours",
			input: DurationString(24 * time.Hour),
			want:  `"24h0m0s"`,
		},
		{
			name:  "Mixed",
			input: DurationString(24*time.Hour + 30*time.Minute + 15*time.Second),
			want:  `"24h30m15s"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Errorf("Duration.MarshalJSON() = %q, want %q", string(got), tt.want)
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    DurationString
		wantErr bool
	}{
		{
			name:  "Hours",
			input: `"4h"`,
			want:  DurationString(4 * time.Hour),
		},
		{
			name:  "Mixed",
			input: `"2h30m"`,
			want:  DurationString(2*time.Hour + 30*time.Minute),
		},
		{
			name:  "Full",
			input: `"6h0m0s"`,
			want:  DurationString(6 * time.Hour),
		},
		{
			name:    "Invalid",
			input:   "daily",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DurationString
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Error = %q, want error: %t", err, tt.wantErr)
			}
			if got.String() != tt.want.String() {
				t.Errorf("Duration.UnmarshalJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}
