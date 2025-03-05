package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jevvonn/readora-backend/config"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"github.com/jevvonn/readora-backend/internal/infra/mailer"
)

// A list of task types.
const (
	SendOTPRegisterTaskName = "email:send-otp-register"
)

type SendOTPRegisterPayload struct {
	Email string
	OTP   string
}

// Create Task
func HandleSendOTPRegisterTask(ctx context.Context, t *asynq.Task) error {
	var payload SendOTPRegisterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log := logger.New()
	conf := config.Load()

	mailer := mailer.New()
	OTPLink := fmt.Sprintf("%s/auth/verify-otp?email=%s&otp=%s", conf.FrontendBaseURL, payload.Email, payload.OTP)
	err := mailer.Send([]string{payload.Email}, "Register OTP Link", fmt.Sprintf(`
		<h3>Register OTP Link</h3>
		<p>Here is your OTP Link: <a href="%s">Click Here</a></p>
		
		<p>Or you can copy this link: %s</p>
	`, OTPLink, OTPLink))

	if err != nil {
		log.Error("[Task][SendOtpRegister]", err)
		return err
	}

	log.Info("[Task][SendOtpRegister]", "Sending Email to User: email="+payload.Email)
	fmt.Println("[Task][SendOtpRegister] Run Succesfully")
	return nil
}
