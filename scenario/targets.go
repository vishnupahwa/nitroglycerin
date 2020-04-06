package scenario

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"sync/atomic"
)

func StaticInterceptedTargeter(override string, interceptor func(tgt vegeta.Target) vegeta.Target, tgts ...vegeta.Target) vegeta.Targeter {
	if override != "" {
		log.Println("Overriding all target URIs with " + override)
		for i, tgt := range tgts {
			tgt.URL = override
			tgts[i] = tgt
		}
	}
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
