package limiter

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

func New(burst int, interval time.Duration, tokens int) *Limiter {
	r := float64(burst) / interval.Seconds()
	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(r), burst),
		tokens:  tokens,
	}
}

type Limiter struct {
	limiter *rate.Limiter
	tokens  int
}

func (l *Limiter) Wait(ctx context.Context) error {
	return l.limiter.WaitN(ctx, l.tokens)
}

func (l *Limiter) Allow() bool {
	return l.limiter.AllowN(time.Now(), l.tokens)
}
