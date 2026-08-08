//go:build kyse

package components

{{-- The theme is client state -- which colours somebody prefers on this device
     -- so the server neither sets it nor is told about it. It lives in
     localStorage and on the html element, which is the case RULE 13 allows
     Alpine for: "estado do cliente, efêmero e invisível ao servidor".

     It takes no props, and still earns its place: the markup is a button, a menu
     of six swatches and the ARIA that ties them together, and repeating that on
     every layout is how two of them drift apart.

     The store it reads is set up by theme.js, which the framework embeds and the
     layout loads before the body is parsed. Nothing here declares x-data for the
     theme itself: theme.js applies it to <html> ahead of the first paint, and
     binding it again in markup names a component that does not exist. --}}

<div class="dropdown-menu" x-data>
	<button type="button" class="btn" data-variant="ghost" data-size="icon"
		aria-label="Change the theme" aria-haspopup="menu" aria-expanded="false">
		<span aria-hidden="true" x-text="$store.theme.dark ? '◑' : '○'"></span>
	</button>

	<div data-popover aria-hidden="true" class="w-44">
		<div role="menu">
			<button type="button" role="menuitem" @click="$store.theme.toggleDark()">
				<span x-text="$store.theme.dark ? 'Light' : 'Dark'"></span>
			</button>

			<hr role="separator">

			<template x-for="a in $store.theme.accents" :key="a">
				<button type="button" role="menuitem" @click="$store.theme.setAccent(a)"
					:aria-current="$store.theme.accent === a">
					<span class="size-3 rounded-full border" :style="`background: var(--accent-swatch-${a})`"></span>
					<span x-text="a"></span>
				</button>
			</template>
		</div>
	</div>
</div>
