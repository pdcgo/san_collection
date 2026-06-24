package san_verification

// MockOtpCode is the one-time code the mock OtpVerification accepts, so local dev
// and tests can exercise the OTP flow without sending real SMS.
const MockOtpCode = "123456"

type mockOtpVerification struct{}

// NewMockOtpVerification returns an OtpVerification that never contacts a real SMS
// provider: Send is a no-op and Verify approves only MockOtpCode.
func NewMockOtpVerification() OtpVerification {
	return &mockOtpVerification{}
}

// Send implements [OtpVerification]. It is a no-op — no SMS is sent.
func (m *mockOtpVerification) Send(phone string) error {
	return nil
}

// Verify implements [OtpVerification]. It approves only MockOtpCode, regardless
// of phone.
func (m *mockOtpVerification) Verify(code, phone string) (bool, error) {
	return code == MockOtpCode, nil
}
