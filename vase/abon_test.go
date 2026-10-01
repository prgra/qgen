package vase

import "testing"

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{name: "international prefix", phone: "+7 (999) 123-45-67", want: "+79991234567"},
		{name: "domestic prefix", phone: "8 (999) 123-45-67", want: "+79991234567"},
		{name: "without prefix", phone: "999 123-45-67", want: "+79991234567"},
		{name: "invalid country code", phone: "+1 999 123-4567", want: ""},
		{name: "invalid length", phone: "999123", want: ""},
		{name: "invalid characters", phone: "tel9991234567", want: ""},
		{name: "empty", phone: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePhone(tt.phone); got != tt.want {
				t.Errorf("normalizePhone(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}
