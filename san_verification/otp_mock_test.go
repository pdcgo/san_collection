package san_verification

import "testing"

func TestMockOtpVerification(t *testing.T) {
	v := NewMockOtpVerification()

	if err := v.Send("+628123456789"); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	ok, err := v.Verify(MockOtpCode, "+628123456789")
	if err != nil {
		t.Fatalf("Verify(correct) returned error: %v", err)
	}
	if !ok {
		t.Fatalf("Verify(%q) = false; want true", MockOtpCode)
	}

	bad, err := v.Verify("000000", "+628123456789")
	if err != nil {
		t.Fatalf("Verify(wrong) returned error: %v", err)
	}
	if bad {
		t.Fatalf("Verify(\"000000\") = true; want false")
	}
}
