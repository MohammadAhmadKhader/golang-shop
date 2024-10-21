package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key       string
	CreatedAt time.Time
}

type RetentionMap map[string]OTP

var mu sync.Mutex

func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap, 0)

	go rm.Retention(ctx, retentionPeriod)
	return rm
}

func (r RetentionMap) NewOTP() OTP {
	otp := OTP{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}

	r[otp.Key] = otp
	return otp
}

func (r RetentionMap) ValidateOTP(OTPkey string) bool {
	_, ok := r[OTPkey]
	if !ok {
		return false
	}

	r.deleteOTP(OTPkey)

	return true
}

func (r RetentionMap) Retention(ctx context.Context, retentionPeriod time.Duration) {
	ticker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			for _, otp := range r {
				if otp.CreatedAt.Add(retentionPeriod).Before(time.Now()) {
					r.deleteOTP(otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r RetentionMap) deleteOTP(OTPkey string) {
	mu.Lock()
	delete(r, OTPkey)
	mu.Unlock()
}
