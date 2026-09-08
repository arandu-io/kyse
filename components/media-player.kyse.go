//go:build kyse

package components

@go
// MediaPlayerProps is a video or a piece of audio, played by the browser's own
// player.
//
// The controls are the browser's, and that is a decision rather than a
// shortcut. A custom player has to rebuild play, scrub, volume, captions,
// fullscreen and picture-in-picture, each with its keyboard contract, and then
// keep them working against a media element whose state changes for reasons
// the page did not cause -- a track that stalls, a call that takes audio
// focus, a headset that pauses it. Every one of those is already right in the
// element.
//
// Captions are a track and not a burned-in subtitle, so they can be turned
// off, restyled, translated and read by anything that is not looking at the
// screen. A video with speech and no track is a video half the audience cannot
// use.
//
// It publishes root, media and caption.
type MediaPlayerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Audio draws an audio element instead of a video one.
	Audio bool
	// Sources are the files, most preferred first. The browser plays the first
	// it can decode, so a modern codec can lead with an older one behind it.
	Sources []MediaSource
	// Tracks are the caption, subtitle and description files.
	Tracks []MediaTrack
	// PosterURL is the still shown before a video plays. Without one the
	// element shows the first frame, which is often black.
	PosterURL string
	// Label names the player. A page with several needs it, and a media
	// element with no name announces as "video".
	Label string
	// Ratio is the shape the video is framed in: "16/9", "4/3", "1/1". Empty
	// leaves the file's own, which makes the page jump when the metadata
	// arrives.
	Ratio string
	// Autoplay starts without being asked. It only works muted in every
	// current browser, so Muted is forced on with it -- and it is worth not
	// using: audio that starts on its own is the single most complained-about
	// behaviour on the web.
	Autoplay bool
	// Loop restarts at the end.
	Loop bool
	// Muted starts silent.
	Muted bool
	// Preload is how much is fetched before play: "none", "metadata" or
	// "auto". Empty fetches the metadata, which is enough for the duration and
	// the scrubber and is a few kilobytes.
	Preload string
	// Caption is the sentence under the player.
	Caption string
}

// MediaSource is one file the browser may play.
type MediaSource struct {
	// URL is the file.
	URL string
	// Type is the MIME type, with the codecs where they matter:
	// `video/mp4; codecs="av01.0.05M.08"`. It is what lets the browser skip a
	// file it cannot play without downloading any of it.
	Type string
}

// MediaTrack is one text track: captions, subtitles or descriptions.
type MediaTrack struct {
	// URL is the WebVTT file.
	URL string
	// Kind is "captions", "subtitles", "descriptions" or "chapters". Empty is
	// captions, which is the one a video with speech needs.
	Kind string
	// Language is the BCP 47 tag: "en", "pt-BR".
	Language string
	// Label is what the track is called in the browser's own menu.
	Label string
	// Default turns this track on without being asked.
	Default bool
}

// Silent is whether the player starts muted. Autoplay forces it, because every
// current browser refuses to autoplay with sound and the refusal is silent --
// the video simply does not start, and nothing says why.
func (p MediaPlayerProps) Silent() bool { return p.Muted || p.Autoplay }

// Fetch is how much is loaded before play.
func (p MediaPlayerProps) Fetch() string {
	if p.Preload != "" {
		return p.Preload
	}
	return "metadata"
}

// TrackKind is what one track is, defaulting to captions.
func (p MediaPlayerProps) TrackKind(track MediaTrack) string {
	if track.Kind != "" {
		return track.Kind
	}
	return "captions"
}

// PartNames are the parts this component publishes.
func (p MediaPlayerProps) PartNames() []string {
	return []string{"root", "media", "caption"}
}
@endgo

<figure
	data-part="root"
	class="{{ .RootClass("media-player") }}"
	@if(.Ratio != "")
		data-ratio="{{ .Ratio }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(.Audio)
		<audio
			data-part="media"
			@if(.PartClass("media") != "")
				class="{{ .PartClass("media") }}"
			@endif
			controls
			preload="{{ .Fetch() }}"
			@attributes(.PartAttrs("media"))
			@if(.Label != "")
				aria-label="{{ .Label }}"
			@endif
			@if(.Autoplay)
				autoplay
			@endif
			@if(.Loop)
				loop
			@endif
			@if(.Silent())
				muted
			@endif
		>
			@foreach(.Sources as source)
				<source
					src="{{ source.URL }}"
					@if(source.Type != "")
						type="{{ source.Type }}"
					@endif
				>
			@endforeach
			@foreach(.Tracks as track)
				<track
					src="{{ track.URL }}"
					kind="{{ .TrackKind(track) }}"
					srclang="{{ track.Language }}"
					label="{{ track.Label }}"
					@if(track.Default)
						default
					@endif
				>
			@endforeach
		</audio>
	@endif
	@if(!.Audio)
		<video
			data-part="media"
			class="{{ .PartClass("media", "media-player-video") }}"
			controls
			playsinline
			preload="{{ .Fetch() }}"
			@attributes(.PartAttrs("media"))
			@if(.PosterURL != "")
				poster="{{ .PosterURL }}"
			@endif
			@if(.Label != "")
				aria-label="{{ .Label }}"
			@endif
			@if(.Autoplay)
				autoplay
			@endif
			@if(.Loop)
				loop
			@endif
			@if(.Silent())
				muted
			@endif
		>
			@foreach(.Sources as source)
				<source
					src="{{ source.URL }}"
					@if(source.Type != "")
						type="{{ source.Type }}"
					@endif
				>
			@endforeach
			@foreach(.Tracks as track)
				<track
					src="{{ track.URL }}"
					kind="{{ .TrackKind(track) }}"
					srclang="{{ track.Language }}"
					label="{{ track.Label }}"
					@if(track.Default)
						default
					@endif
				>
			@endforeach
		</video>
	@endif

	@if(.Caption != "")
		<figcaption
			data-part="caption"
			class="{{ .PartClass("caption", "figure-caption") }}"
			@attributes(.PartAttrs("caption"))
		>{{ .Caption }}</figcaption>
	@endif
</figure>
