package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Profiling is enabled by setting MATCHA_PROFILE_DIR. It writes:
//   - cpu.pprof: a CPU profile for the whole session
//   - slow.log:  every Update/View call that took longer than one frame
const slowThreshold = 16 * time.Millisecond

type profiledModel struct {
	inner tea.Model
	log   *os.File
}

func startProfiling(m tea.Model) (tea.Model, func()) {
	dir := os.Getenv("MATCHA_PROFILE_DIR")
	if dir == "" {
		return m, func() {}
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "profiling disabled: %v\n", err)
		return m, func() {}
	}
	cpu, err := os.Create(filepath.Join(dir, "cpu.pprof"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "profiling disabled: %v\n", err)
		return m, func() {}
	}
	if err := pprof.StartCPUProfile(cpu); err != nil {
		fmt.Fprintf(os.Stderr, "profiling disabled: %v\n", err)
		_ = cpu.Close()
		return m, func() {}
	}
	slow, err := os.Create(filepath.Join(dir, "slow.log"))
	if err != nil {
		slow = nil
	}
	return &profiledModel{inner: m, log: slow}, func() {
		pprof.StopCPUProfile()
		_ = cpu.Close()
		if slow != nil {
			_ = slow.Close()
		}
	}
}

func (p *profiledModel) record(kind string, msg tea.Msg, d time.Duration) {
	if p.log == nil || d < slowThreshold {
		return
	}
	detail := ""
	if msg != nil {
		detail = fmt.Sprintf("%T", msg)
	}
	if k, ok := msg.(tea.KeyPressMsg); ok {
		detail += " " + k.String()
	}
	_, _ = fmt.Fprintf(p.log, "%s %-6s %8s %s screen=%T\n",
		time.Now().Format("15:04:05.000"), kind, d.Round(time.Microsecond), detail, p.screen())
}

func (p *profiledModel) screen() any {
	if s, ok := p.inner.(interface{ Current() tea.Model }); ok {
		return s.Current()
	}
	return p.inner
}

func (p *profiledModel) Init() tea.Cmd { return p.inner.Init() }

func (p *profiledModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	start := time.Now()
	var cmd tea.Cmd
	p.inner, cmd = p.inner.Update(msg)
	p.record("update", msg, time.Since(start))
	return p, cmd
}

func (p *profiledModel) View() tea.View {
	start := time.Now()
	v := p.inner.View()
	p.record("view", nil, time.Since(start))
	return v
}
