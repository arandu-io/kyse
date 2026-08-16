// Package mailui draws the messages an application sends.
//
// It is the set every project rewrites otherwise, and rewrites badly: a mail
// client is a browser from 2005 with the interesting parts removed, and the
// markup that survives it is table-based, inline-styled and nothing like the
// page these same components draw on the web.
//
//	appmail.Content{View: "mail.verify", Data: m}
//
//	{!! mailui.Layout(mailui.LayoutProps{
//	    Brand:   "Arandu",
//	    Heading: "Confirm your email address",
//	    Body: mailui.Paragraph("Hello Ada, and welcome.") +
//	        mailui.Button(mailui.ButtonProps{Label: "Confirm my address", Href: link}) +
//	        mailui.Small("The link works for 24 hours."),
//	    Footer: "If you did not create this account, ignore this message.",
//	}) !!}
//
// # Why this is not the component package
//
// components/ draws a page: semantic classes, a stylesheet, a cascade. None of
// that reaches an e-mail. Outlook renders with Word, Gmail strips <style>, and
// an external stylesheet is not fetched at all -- so every rule here is an
// inline attribute on the element it applies to, and the layout is a table
// because a table is the one box model every client agrees on.
//
// Two packages rather than a mode on one, because the constraint is not a
// preference: a component that tried to serve both would carry the 2005 markup
// into every page on the web, or the modern markup into every inbox.
//
// # What it deliberately does not do
//
// No dark mode. `prefers-color-scheme` reaches three clients out of a dozen, and
// a message that is half inverted is worse than one that is not. The palette is
// light and states so.
//
// No images. Every client blocks remote images until the reader allows them, so
// a logo delivered that way is a broken-image box on the first impression --
// which, for a verification mail, is the whole job. The brand is set in type.
//
// No markdown. A view already produces both the HTML and the text part, so a
// second templating language for e-mail would be a second way to draw the same
// message.
package mailui

import (
	"fmt"
	"html/template"
	"strings"
)

// The palette, and it is inline everywhere because it has to be.
//
// Named here so that a message and the next message agree, which is the whole
// difference between an application that looks like one thing and a set of
// e-mails that each look like whoever wrote them.
const (
	ink     = "#18181b"
	body    = "#3f3f46"
	muted   = "#71717a"
	faint   = "#a1a1aa"
	rule    = "#e4e4e7"
	surface = "#ffffff"
	page    = "#f4f4f5"

	// font is a stack and not a face. A vendored woff2 does not reach an inbox:
	// @font-face is stripped by most clients and fetched by none of the ones
	// that keep it, so the type is whatever the reader's machine has.
	font = `-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Helvetica,Arial,sans-serif`
)

// LayoutProps is the message around the content.
type LayoutProps struct {
	// Brand is the name at the top, set in type rather than fetched as an image.
	Brand string
	// Heading is the first line inside the card, and it is what the message is
	// about. A message with no heading is one the reader has to parse a
	// paragraph to understand.
	Heading string
	// Body is what the helpers below produced.
	Body template.HTML
	// Footer is the grey line under the card: why this arrived, and what to do
	// if it should not have. Every message an application sends needs one, and
	// its absence is what a spam complaint is made of.
	Footer string
	// Preheader is the sentence a client shows in the inbox list beside the
	// subject. Left empty, it shows the first words of the body -- usually
	// "Hello there," which tells nobody anything.
	Preheader string
}

// Layout renders the whole message.
//
// A table, because a table is the one box model every client agrees on: Outlook
// renders with Word, which has no flexbox, no grid and an idea of margin from
// another decade. The 560px width is what fits a phone at full bleed and a
// desktop client's reading pane without a horizontal scroll.
func Layout(p LayoutProps) template.HTML {
	var b strings.Builder

	fmt.Fprintf(&b, `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="light">
<title>%s</title>
</head>
<body style="margin:0;padding:24px;background:%s;font-family:%s;color:%s;-webkit-font-smoothing:antialiased">
`, escape(p.Heading), page, font, ink)

	if p.Preheader != "" {
		// Hidden in the message and read by the inbox list. The trailing
		// entities stop the client from filling the rest of the preview line
		// with whatever comes next in the markup.
		fmt.Fprintf(&b, `<div style="display:none;max-height:0;overflow:hidden;opacity:0">%s%s</div>
`, escape(p.Preheader), strings.Repeat("&#847;&zwnj;&nbsp;", 30))
	}

	fmt.Fprintf(&b, `<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
<tr><td align="center">
<table role="presentation" width="560" cellpadding="0" cellspacing="0" border="0" style="width:100%%;max-width:560px">
`)

	if p.Brand != "" {
		fmt.Fprintf(&b, `<tr><td style="padding:0 4px 20px;font-size:15px;font-weight:600;letter-spacing:-0.01em;color:%s">%s</td></tr>
`, ink, escape(p.Brand))
	}

	fmt.Fprintf(&b, `<tr><td style="background:%s;border:1px solid %s;border-radius:10px;padding:32px">
`, surface, rule)

	if p.Heading != "" {
		fmt.Fprintf(&b, `<h1 style="margin:0 0 16px;font-size:20px;line-height:1.3;font-weight:600;letter-spacing:-0.02em;color:%s">%s</h1>
`, ink, escape(p.Heading))
	}

	b.WriteString(string(p.Body))
	b.WriteString("</td></tr>\n")

	if p.Footer != "" {
		fmt.Fprintf(&b, `<tr><td style="padding:16px 8px;font-size:12px;line-height:1.6;color:%s">%s</td></tr>
`, faint, escape(p.Footer))
	}

	b.WriteString("</table>\n</td></tr>\n</table>\n</body>\n</html>\n")
	return template.HTML(b.String())
}

// Paragraph is one block of text.
func Paragraph(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<p style="margin:0 0 16px;font-size:15px;line-height:1.6;color:%s">%s</p>`,
		body, escape(text)))
}

// ButtonProps is the one thing the message wants the reader to do.
type ButtonProps struct {
	Label string
	Href  string
}

// Button is the call to action.
//
// One per message. A message with two buttons has no primary action, and the
// reader picks the wrong one or neither -- which is why the plain-text part
// carries one link and not a list.
//
// It is an <a> styled as a block and not a <button>: a button element does
// nothing in an inbox, and a form does not submit from one.
func Button(p ButtonProps) template.HTML {
	if p.Href == "" || p.Label == "" {
		return ""
	}
	return template.HTML(fmt.Sprintf(
		`<p style="margin:0 0 24px"><a href="%s" style="display:inline-block;background:%s;color:#ffffff;text-decoration:none;padding:12px 20px;border-radius:8px;font-size:15px;font-weight:500;line-height:1">%s</a></p>`,
		escape(p.Href), ink, escape(p.Label)))
}

// Fallback is the address under a button, written out.
//
// It is not redundancy. A button is an <a> with a background, and a client that
// strips the style leaves a bare word; one that blocks the link leaves nothing.
// The URL in text is what the reader copies when either happens, and it is the
// difference between a verification that can be completed and one that cannot.
func Fallback(href string) template.HTML {
	if href == "" {
		return ""
	}
	return template.HTML(fmt.Sprintf(
		`<p style="margin:0;font-size:13px;line-height:1.6;color:%s">If the button does not open, copy this address:<br><span style="word-break:break-all">%s</span></p>`,
		muted, escape(href)))
}

// Small is the quiet line: an expiry, a note, a condition.
func Small(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<p style="margin:0 0 16px;font-size:13px;line-height:1.6;color:%s">%s</p>`,
		muted, escape(text)))
}

// Panel is a block set apart: a code, a quoted line, a summary.
func Panel(text string) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 16px"><tr><td style="background:%s;border-radius:8px;padding:16px;font-size:15px;line-height:1.5;color:%s">%s</td></tr></table>`,
		page, ink, escape(text)))
}

// Divider is a rule between two parts of a message.
func Divider() template.HTML {
	return template.HTML(fmt.Sprintf(
		`<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 16px"><tr><td style="border-top:1px solid %s;font-size:0;line-height:0">&nbsp;</td></tr></table>`,
		rule))
}

// escape is html/template's escaping, applied by hand.
//
// Everything here returns template.HTML, which tells the engine the value is
// already safe -- so the engine escapes nothing, and every string reaching one
// of these functions has to be escaped here. A Label or a Heading is often a
// person's name or a subject somebody typed.
func escape(s string) string { return template.HTMLEscapeString(s) }
