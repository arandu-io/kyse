//go:build kyse

package components

@go
// ToastProps is what just happened, said once and then gone.
//
// It is written to be returned on its own, not only drawn inside a page. An
// endpoint that saves something answers with this fragment and
//
//	hx-target="#toaster" hx-swap="beforeend"
//
// appends it to the tray the layout already has. That is the whole flash-message
// mechanism: no session key, no redirect carrying a string, no template deciding
// whether to draw an empty box.
//
// The behaviour comes from the vendored script rather than from anything here.
// It arms every .toast it finds, dismisses after data-duration, pauses the timer
// while the pointer is over the tray and removes the node when the transition
// ends. A toast written by hand gets one of those four right.
type ToastProps struct {
	// Title is the line. Empty draws nothing at all, so a page can include it
	// unconditionally.
	Title string
	// Message is the detail under it, and is usually not needed.
	Message string

	// Category is "error" for a failure and empty for anything else. It picks
	// both the styling and the default duration -- a failure stays five seconds,
	// a confirmation three.
	Category string
	// Duration overrides that, in milliseconds. Zero takes the default; a
	// negative number stays until dismissed, which is right for anything that
	// needs an answer.
	Duration int

	// ActionLabel and ActionURL draw a button inside the toast: undo, retry,
	// view. Both empty draws none.
	ActionLabel string
	ActionURL   string
}

// Live is how insistently assistive technology announces it. A failure
// interrupts; a confirmation waits its turn.
func (p ToastProps) Live() string {
	if p.Category == "error" {
		return "assertive"
	}
	return "polite"
}
@endgo

@if(.Title != "")
	<div class="toast" role="status" aria-live="{{ .Live() }}"
		@if(.Category != "")
			data-category="{{ .Category }}"
		@endif
		@if(.Duration != 0)
			data-duration="{{ .Duration }}"
		@endif
	>
		<div class="toast-content">
			<section>
				<h2>{{ .Title }}</h2>
				@if(.Message != "")
					<p>{{ .Message }}</p>
				@endif
			</section>

			@if(.ActionURL != "")
				<footer>
					<button type="button" class="btn" data-size="sm" data-toast-action
						hx-post="{{ .ActionURL }}" hx-swap="none">{{ .ActionLabel }}</button>
				</footer>
			@endif
		</div>
	</div>
@endif
