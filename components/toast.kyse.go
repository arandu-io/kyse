//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

@go
// ToastProps is what just happened, said once and then gone.
//
// It reaches the page by the road the action took.
//
// An action that redirects -- a form that worked, which is answered with a
// redirect like every form -- says what it did with Context.With("status",
// ...), and the page it lands on carries the sentence in its flash. The layout
// draws it with one call, on every page, because a page with nothing flashed
// draws nothing:
//
//	{!! components.Toast(components.ToastProps{Page: .Page}) !!}
//
// An endpoint that answers in place, without navigating, returns this fragment
// instead, and
//
//	hx-target="#toaster" hx-swap="beforeend"
//
// appends it to the tray the layout already has. The two are one toast reached
// from the two answers an action can give; neither needs a template deciding
// whether to draw an empty box.
//
// The behaviour comes from the vendored script rather than from anything here.
// It arms every .toast it finds, dismisses after data-duration, pauses the timer
// while the pointer is over the tray and removes the node when the transition
// ends. A toast written by hand gets one of those four right.
type ToastProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Title is the line. Empty takes what Page flashed; with nothing flashed
	// either it draws nothing at all, so a page can include it
	// unconditionally.
	Title string
	// Message is the detail under it, and is usually not needed.
	Message string

	// Page is the screen's own view.Page. With Title empty the toast draws
	// the sentence the request before this one flashed under FlashKey -- what
	// Context.With left -- and nothing when it flashed none. The close button
	// reads its word from it too, when CloseLabel is empty.
	Page Page
	// FlashKey is the key that sentence was flashed under. Empty means
	// "status", the key a form that worked says what it did under.
	FlashKey string

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
	// own. Empty says kyse.toast.close from the Page, or "Close".
	CloseLabel string
}

// Headline is the line the toast draws: Title, or else what the page flashed
// under FlashKey. Empty draws no toast.
func (p ToastProps) Headline() string {
	if p.Title != "" {
		return p.Title
	}
	flashes, ok := p.Page.(FlashPage)
	if !ok {
		return ""
	}
	key := p.FlashKey
	if key == "" {
		key = "status"
	}
	return flashes.Flashed(key)
}

// Close is the accessible name of the close button.
func (p ToastProps) Close() string {
	if p.CloseLabel == "" {
		return say(p.Page, "toast.close", "Close")
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

@if(.Headline() != "")
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
				>{{ .Headline() }}</h2>
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
