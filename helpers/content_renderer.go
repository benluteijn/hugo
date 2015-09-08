package helpers

import (
	"bytes"
	"html"

	"github.com/russross/blackfriday"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type LinkResolverFunc func(ref string) (string, error)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
type HugoHtmlRenderer struct {
	LinkResolver LinkResolverFunc
	blackfriday.Renderer
}

func (renderer *HugoHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if viper.GetBool("PygmentsCodeFences") {
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, ""))
	} else {
		renderer.Renderer.BlockCode(out, text, lang)
	}
}

func (renderer *HugoHtmlRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if renderer.LinkResolver == nil || bytes.HasPrefix(link, []byte("{@{@HUGOSHORTCODE")) {
		// Use the blackfriday built in Links
		renderer.Renderer.Link(out, link, title, content)
	} else {
		newLink, err := renderer.LinkResolver(string(link))
		if err != nil {
			newLink = string(link)
			jww.ERROR.Printf("GH: %s", err)
		}

		renderer.Renderer.Link(out, []byte(newLink), title, content)
	}
}

func (renderer *HugoHtmlRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	renderer.Renderer.Image(out, link, title, alt)
}
