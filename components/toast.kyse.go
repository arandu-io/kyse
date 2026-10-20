//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

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
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
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

	// Dismissible draws a close button in the corner, so the person can send
	// the toast away before its timer does. The vendored script closes a toast
	// on any button in its footer, so the button needs nothing of its own.
	Dismissible bool
	// CloseLabel names that button, which is an icon and has no text of its
	// own. Empty means "Close".
	CloseLabel string
}

// Close is the accessible name of the close button.
func (p ToastProps) Close() string {
	if p.CloseLabel == "" {
		return "Close"
	}
	return p.CloseLabel
}

// Live is how insistently assistive technology announces it. A failure
// interrupts; a confirmation waits its turn.
func (p ToastProps) Live() string {
	if p.Category == "error" {
		return "assertive"
	}
	return "polite"
}
// PartNames are the parts this component publishes.
func (p ToastProps) PartNames() []string {
	return []string{"root", "content", "title", "message", "action", "close"}
}
@endgo

@if(.Title != "")
	<div
		data-part="root"
		class="{{ .RootClass("toast") }}"
		role="status"
		aria-live="{{ .Live() }}"
		@if(.Category != "")
			data-category="{{ .Category }}"
		@endif
		@if(.Duration != 0)
			data-duration="{{ .Duration }}"
		@endif
		@attributes(.RootAttrs())
	>
		<div
			data-part="content"
			class="{{ .PartClass("content", "toast-content") }}"
			@attributes(.PartAttrs("content"))
		>
			<section>
				<h2
					data-part="title"
					@if(.PartClass("title") != "")
						class="{{ .PartClass("title") }}"
					@endif
					@attributes(.PartAttrs("title"))
				>{{ .Title }}</h2>
				@if(.Message != "")
					<p
						data-part="message"
						@if(.PartClass("message") != "")
							class="{{ .PartClass("message") }}"
						@endif
						@attributes(.PartAttrs("message"))
					>{{ .Message }}</p>
				@endif
			</section>

			@if(.ActionURL != "" || .Dismissible)
				<footer>
					@if(.ActionURL != "")
					<button
						data-part="action"
						class="{{ .PartClass("action", "btn") }}"
						type="button"
						data-size="sm"
						data-toast-action
						hx-post="{{ .ActionURL }}"
						hx-swap="none"
						@attributes(.PartAttrs("action"))
					>{{ .ActionLabel }}</button>
					@endif
					@if(.Dismissible)
						<button
							data-part="close"
							class="{{ .PartClass("close", "btn") }}"
							type="button"
							data-variant="ghost"
							data-size="icon-sm"
							data-toast-cancel
							aria-label="{{ .Close() }}"
							@attributes(.PartAttrs("close"))
						>{!! icons.X(icons.Props{}) !!}</button>
					@endif
				</footer>
			@endif
		</div>
	</div>
@endif
