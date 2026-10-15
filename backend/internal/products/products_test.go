package products

import "testing"

func TestFormatCode(t *testing.T) {
	got := FormatCode("ABCDE12345ABCDE12345ABCDE12345")
	want := "ABCDE-12345-ABCDE-12345-ABCDE-12345"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
func TestValidCode(t *testing.T) {
	cases := []struct {
		code string
		want bool
	}{{"ABCDE12345ABCDE12345ABCDE12345", true}, {"abcde12345ABCDE12345ABCDE12345", false}, {"ABCDE-12345-ABCDE-12345-ABCDE-12345", false}, {"ABCDE12345", false}, {"ABCDE12345ABCDE12345ABCDE1234!", false}}
	for _, tc := range cases {
		if got := validCode(tc.code); got != tc.want {
			t.Errorf("validCode(%q)=%v want %v", tc.code, got, tc.want)
		}
	}
}
