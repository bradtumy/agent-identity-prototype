package policy

import "testing"

func TestIsActionAllowedForRole(t *testing.T) {
	if !IsActionAllowedForRole("data-fetcher", "fetch_data") {
		t.Errorf("data-fetcher should be allowed fetch_data")
	}
	if IsActionAllowedForRole("data-fetcher", "notify") {
		t.Errorf("data-fetcher should not be allowed notify")
	}
}

func TestValidatePolicy(t *testing.T) {
	tests := []struct {
		action  string
		role    string
		wantErr bool
	}{
		{"fetch_data", "data-fetcher", false},
		{"fetch_data", "transformer", true},
		{"transform", "transformer", false},
		{"notify", "unknown", true},
		{"unknown", "data-fetcher", true},
	}
	for _, tc := range tests {
		err := ValidatePolicy(tc.action, tc.role)
		if tc.wantErr && err == nil {
			t.Errorf("expected error for action %s role %s", tc.action, tc.role)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("unexpected error for action %s role %s: %v", tc.action, tc.role, err)
		}
	}
}
