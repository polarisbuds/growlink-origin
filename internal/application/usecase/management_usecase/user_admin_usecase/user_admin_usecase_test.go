package user_admin_usecase

import "testing"

func TestValidateLineWorksLocalID(t *testing.T) {
	if err := ValidateLineWorksLocalID("user.name"); err != nil {
		t.Fatalf("expected valid id: %v", err)
	}
	if err := ValidateLineWorksLocalID("bad@id"); err == nil {
		t.Fatal("expected error for @ in id")
	}
}
