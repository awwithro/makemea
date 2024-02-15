package randomtable

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)


func NewMarkdownParser(tree Tree) goldmark.Markdown{
	return goldmark.New(
			goldmark.WithExtensions(extension.GFM,extension.DefinitionList),
			goldmark.WithRendererOptions(
				renderer.WithNodeRenderers(
					util.Prioritized(NewRandomTableRenderer(tree), 1))),
		)
}