package helpers

import (
	"bytes"
	"html"

	"github.com/russross/blackfriday"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type RefLinkFunc func(ref string) (string, error)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
type HugoHtmlRenderer struct {
	Dir     string
	RefLink RefLinkFunc
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
	if renderer.RefLink == nil {
		// Use the blackfriday built in Links
		renderer.Renderer.Link(out, link, title, content)
	} else {
		jww.ERROR.Printf("Sven was rendering a link here (%v, %v, %v)\n", string(link), string(title), renderer.Dir)
		//	renderer.Renderer.Link(out, link, []byte("SVEN"+string(title)), []byte(renderer.Page.Node.Site.RelRef(link, renderer.Page)+string(content)))

		newTitle := []byte("found relref" + string(title))
		newLink, err := renderer.RefLink(string(link))
		if err != nil {
			newLink = string(link)
			newTitle = title
		}

		renderer.Renderer.Link(out, []byte(newLink), newTitle, content)
	}
}
