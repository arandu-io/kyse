package components

import "sync/atomic"

// Config is what every component in a project takes unless it says otherwise.
type Config struct {
	// Theme names a variation the stylesheet defines. It is written on each
	// component as data-theme, and a component given one of its own uses that
	// instead.
	Theme string
}

// Configure sets the defaults for this process. It is called once, from the
// application's boot, before anything is served:
//
//	components.Configure(components.Config{Theme: "admin"})
//
// # Why it panics on the second call
//
// A default that can change while requests are in flight is a page that renders
// one way and, moments later, another -- with nothing recording which. Worse,
// two calls means two places that decide, and reading the code tells you
// neither.
//
// So the second call is a panic rather than an overwrite, the same answer
// view.RegisterStylesheet gives to a second stylesheet. A process may be
// configured, or it may take the defaults; it may not change its mind.
//
// A test that needs a theme passes it on the component rather than configuring
// the process, which is the same thing this makes every caller do.
func Configure(c Config) {
	if !config.CompareAndSwap(nil, &c) {
		panic("components: Configure was called twice. " +
			"The defaults are set once, from the application's boot, so that what a page renders with " +
			"does not depend on when the page was rendered")
	}
}

// configured is the defaults in force, or the zero Config.
//
// An atomic pointer rather than a mutex: this is read on every element of every
// component and written at most once, and a lock on that path would be
// contention bought for a value that never changes. A plain variable would be a
// data race between the boot goroutine and the first request, which
// `go test -race` reports and which is real -- Go gives no ordering between the
// two on its own.
func configured() Config {
	if c := config.Load(); c != nil {
		return *c
	}
	return Config{}
}

var config atomic.Pointer[Config]
