package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
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

	mailer := mailer.New()
	err := mailer.Send([]string{payload.Email}, "Register OTP", "Your OTP is "+payload.OTP)

	if err != nil {
		fmt.Println("[Task][SendOtpRegister] Abort" + err.Error())
		log.Error("[Task][SendOtpRegister]", err)
		return err
	}

	log.Info("[Task][SendOtpRegister]", "Sending Email to User: email="+payload.Email)
	fmt.Println("[Task][SendOtpRegister] Run Succesfully")
	return nil
}
