package san_verification

type OtpVerification interface {
	Send(phone string) error
	Verify(code, phone string) (bool, error)
}
