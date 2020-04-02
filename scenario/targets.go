package scenario

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"sync/atomic"
)

func StaticInterceptedTargeter(interceptor func(tgt vegeta.Target) vegeta.Target, tgts ...vegeta.Target) vegeta.Targeter {
	if interceptor == nil {
		return vegeta.NewStaticTargeter(tgts...)
	}
	i := int64(-1)
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}
		*tgt = interceptor(tgts[atomic.AddInt64(&i, 1)%int64(len(tgts))])
		return nil
	}
}
