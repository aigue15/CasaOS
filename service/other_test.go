package service

import (
	"testing"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"go.uber.org/goleak"
)

func TestSearch(t *testing.T) {
	logger.LogInitConsoleOnly()
	// ecache (imported by CasaOS-Common) starts a package-level cleanup goroutine in init.
	goleak.VerifyNone(t, goleak.IgnoreAnyFunction("github.com/orca-zhang/ecache.init.0.func1"))

	if d, e := NewOtherService().Search("test"); e != nil || d == nil {
		t.Error("then test search error", e)
	}
}
