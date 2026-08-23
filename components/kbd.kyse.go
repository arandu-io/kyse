//go:build kyse

package components

@go
// KbdProps is a key, or the keys of a shortcut held together.
//
// It takes the keys and draws what goes between them. That division is the
// whole component: a caller that supplied the separators would be a caller
// writing markup, and the separator is the part that drifts -- a plus here, a
// space there, a hyphen on the page somebody wrote last.
type KbdProps struct {
	// Keys are the keys, in the order they are held down: {"Ctrl", "Shift",
	// "P"}, or {"⌘", "K"}. One key draws no separator, and none draws nothing
	// at all, so a screen can include the shortcut it may not have.
	Keys []string
}

// keyNames are the keys drawn as a symbol, and what each one is called.
//
// A cap reading ⌘ is announced as whatever that character is named in the voice
// speaking, which is rarely "Command" and is sometimes silence. The name goes
// on the element so what is announced is the key.
var keyNames = map[string]string{
	"⌘": "Command",
	"⇧": "Shift",
	"⌥": "Option",
	"⌃": "Control",
	"⏎": "Enter",
	"↩": "Return",
	"⌫": "Backspace",
	"⌦": "Delete",
	"⎋": "Escape",
	"⇥": "Tab",
	"␣": "Space",
	"↑": "Up arrow",
	"↓": "Down arrow",
	"←": "Left arrow",
	"→": "Right arrow",
}

// Name is what assistive technology says for the key at position i, and empty
// for a key written as a word, which already reads as itself.
func (p KbdProps) Name(i int) string {
	if i < 0 || i >= len(p.Keys) {
		return ""
	}
	return keyNames[p.Keys[i]]
}
@endgo

@if(len(.Keys) > 0)
	<span class="inline-flex items-center gap-1">
		@for(i := 0; i < len(.Keys); i++)
			@if(i > 0)
				<span>+</span>
			@endif
			<kbd
				class="kbd"
				@if(.Name(i) != "")
					aria-label="{{ .Name(i) }}"
				@endif
			>{{ .Keys[i] }}</kbd>
		@endfor
	</span>
@endif
