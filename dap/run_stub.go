//go:build !fullruntime

package dap

import "fmt"

func (s *session) beginDebugRun() {
	s.sendTerminated(fmt.Errorf("debug run requires moonrun (full runtime build); use moonrun --dap instead of moonbasic --dap"))
}
