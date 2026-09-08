//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/kyse/icons"
)

@go
// TreeBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const TreeBehavior = "tree"

// TreeProps is a hierarchy that opens and closes: a file browser, a category
// list, an organisation chart.
//
// The whole tree is one tab stop, and the arrow keys move through it -- down
// and up between the visible rows, right to open a branch and then to step
// into it, left to close it and then to step out to the parent. Home and End
// go to the ends. That contract is the reason it is a tree and not a list of
// nested details: a hundred rows as a hundred tab stops is a hundred presses
// to get past the widget.
//
// The rows are flat in the markup and their depth is an attribute, not a
// nesting. That is what lets the keyboard move from the last child of one
// branch to the next branch without walking back out of three containers, and
// it is what lets a row be swapped in on its own.
//
// It publishes root, item, toggle, icon and label.
type TreeProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what every row hangs its id off. Two trees on one page need two.
	ID string
	// Label names the tree. Without it it announces with no subject.
	Label string
	// Nodes are the top-level rows, each carrying its own children.
	Nodes []TreeNode
	// Multiple lets more than one row be selected at a time.
	Multiple bool
}

// TreeNode is one row, and what is under it.
type TreeNode struct {
	// Label is the text of the row.
	Label string
	// URL makes the row a link. Empty makes it a row that only selects.
	URL string
	// Icon is drawn before the label: a folder, a file, a kind.
	Icon template.HTML
	// Children are the rows under it. None makes it a leaf, and a leaf has no
	// expanded state at all -- which is what aria-expanded being absent means,
	// and is different from it being false.
	Children []TreeNode
	// Expanded is whether a branch starts open. It means nothing on a leaf.
	Expanded bool
	// Selected marks the row as chosen.
	Selected bool
	// Disabled draws it unavailable, and the arrow keys pass over it.
	Disabled bool
}

// TreeRow is one node flattened, with everything the markup needs to say where
// it sits.
type TreeRow struct {
	// Node is the row itself.
	Node TreeNode
	// Level is how deep it is, counting from one -- which is what aria-level
	// counts from.
	Level int
	// Position is which child it is among its siblings, counting from one.
	Position int
	// Size is how many siblings it has, including itself.
	Size int
	// Index is its place in the flattened list, used only to build an id.
	Index int
	// Branch is whether it has children.
	Branch bool
	// Hidden is whether an ancestor of it is closed. A row under a closed
	// branch is not drawn at all rather than drawn and hidden: aria-expanded
	// on the parent promises the children are not in the tree, and a screen
	// reader that finds them anyway reads a tree that disagrees with itself.
	Hidden bool
}

// Rows is the tree flattened, in the order somebody reading it moves through.
func (p TreeProps) Rows() []TreeRow {
	rows := make([]TreeRow, 0, len(p.Nodes))
	var walk func(nodes []TreeNode, level int, visible bool)
	walk = func(nodes []TreeNode, level int, visible bool) {
		for at, node := range nodes {
			branch := len(node.Children) > 0
			rows = append(rows, TreeRow{
				Node:     node,
				Level:    level,
				Position: at + 1,
				Size:     len(nodes),
				Index:    len(rows),
				Branch:   branch,
				Hidden:   !visible,
			})
			if branch {
				walk(node.Children, level+1, visible && node.Expanded)
			}
		}
	}
	walk(p.Nodes, 1, true)
	return rows
}

// RowID is the id of one row.
func (p TreeProps) RowID(row TreeRow) string {
	return p.ID + "-row-" + strconv.Itoa(row.Index)
}

// Stop is whether this row holds the tab stop: the first visible row that can
// take it.
func (p TreeProps) Stop(row TreeRow) bool {
	for _, one := range p.Rows() {
		if one.Hidden || one.Node.Disabled {
			continue
		}
		return one.Index == row.Index
	}
	return false
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p TreeProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: TreeBehavior, Props: map[string]any{"multiple": p.Multiple}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p TreeProps) PartNames() []string {
	return []string{"root", "item", "toggle", "icon", "label"}
}
@endgo

<ul
	data-part="root"
	class="{{ .RootClass("tree") }}"
	id="{{ .ID }}"
	role="tree"
	aria-label="{{ .Label }}"
	@if(.Multiple)
		aria-multiselectable="true"
	@endif
	@attributes(.RootAttrs())
>
	@foreach(.Rows() as row)
		@if(!row.Hidden)
			<li
				data-part="item"
				class="{{ .PartClass("item", "tree-item") }}"
				id="{{ .RowID(row) }}"
				role="treeitem"
				aria-level="{{ row.Level }}"
				aria-posinset="{{ row.Position }}"
				aria-setsize="{{ row.Size }}"
				data-level="{{ row.Level }}"
				@attributes(.PartAttrs("item"))
				@if(row.Branch && row.Node.Expanded)
					aria-expanded="true"
				@endif
				@if(row.Branch && !row.Node.Expanded)
					aria-expanded="false"
				@endif
				@if(row.Node.Selected)
					aria-selected="true"
				@endif
				@if(row.Node.Disabled)
					aria-disabled="true"
				@endif
				@if(.Stop(row))
					tabindex="0"
				@endif
				@if(!.Stop(row))
					tabindex="-1"
				@endif
			>
				@if(row.Branch)
					<span
						data-part="toggle"
						class="{{ .PartClass("toggle", "tree-toggle") }}"
						aria-hidden="true"
						@attributes(.PartAttrs("toggle"))
					>{!! icons.CaretRight(icons.Props{}) !!}</span>
				@endif
				@if(row.Node.Icon != "")
					<span
						data-part="icon"
						class="{{ .PartClass("icon", "tree-icon") }}"
						aria-hidden="true"
						@attributes(.PartAttrs("icon"))
					>{!! row.Node.Icon !!}</span>
				@endif
				@if(row.Node.URL != "")
					<a
						data-part="label"
						@if(.PartClass("label") != "")
							class="{{ .PartClass("label") }}"
						@endif
						@attributes(.PartAttrs("label"))
						href="{{ row.Node.URL }}"
						tabindex="-1"
					>{{ row.Node.Label }}</a>
				@endif
				@if(row.Node.URL == "")
					<span
						data-part="label"
						@if(.PartClass("label") != "")
							class="{{ .PartClass("label") }}"
						@endif
						@attributes(.PartAttrs("label"))
					>{{ row.Node.Label }}</span>
				@endif
			</li>
		@endif
	@endforeach
</ul>
