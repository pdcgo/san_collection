package san_verification

import (
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioConfiguration struct {
	AppId     string `env:"TWILIO_APP_ID" yaml:"app_id"`
	Token     string `env:"TWILIO_TOKEN" yaml:"token"`
	ServiceId string `env:"TWILIO_SERVICE_ID" yaml:"service_id"`
}

type twilioVerification struct {
	cfg *TwilioConfiguration
}

func NewTwilioOtpVerification(
	cfg *TwilioConfiguration,
) OtpVerification {
	return &twilioVerification{cfg}
}

// Verify implements [OtpVerification].
func (t *twilioVerification) Verify(code string, phone string) (bool, error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: t.cfg.AppId,
		Password: t.cfg.Token,
	})

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(t.cfg.ServiceId,
		params)
	if err != nil {
		return false, err
	}

	if resp.Status == nil {
		return false, nil
	}

	if *resp.Status == "approved" {
		return true, nil
	}

	return false, nil
}

// Send implements [OtpVerification].
func (t *twilioVerification) Send(phone string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: t.cfg.AppId,
		Password: t.cfg.Token,
	})

	params := &verify.CreateVerificationParams{}
	params.SetChannel("sms")
	params.SetTo(phone)

	_, err := client.VerifyV2.CreateVerification(t.cfg.ServiceId,
		params)
	if err != nil {
		return err
	}

	return nil
}
