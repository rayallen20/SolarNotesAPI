package articlecontent

import (
	"bytes"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	htmltoken "golang.org/x/net/html"
)

// Only raw HTML recognized by the Markdown parser can contain HTML image elements.
func (p *imageParser) rewriteHTML(document ast.Node, source []byte) error {
	var literalElement string
	return ast.Walk(document, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		var segments []text.Segment
		switch n := node.(type) {
		case *ast.Image:
			// An image description becomes alt text rather than rendered HTML.
			return ast.WalkSkipChildren, nil
		case *ast.RawHTML:
			segments = n.Segments.Sliced(0, n.Segments.Len())
		case *ast.HTMLBlock:
			lines := n.Lines()
			segments = lines.Sliced(0, lines.Len())
			if n.HasClosure() {
				segments = append(append([]text.Segment(nil), segments...), n.ClosureLine)
			}
		default:
			return ast.WalkContinue, nil
		}
		return ast.WalkContinue, p.rewriteHTMLFragment(segments, source, &literalElement)
	})
}

// Tokenize a contiguous HTML fragment while retaining offsets into the Markdown source.
// Offsets can jump between lines inside a block quote or a list item.
func (p *imageParser) rewriteHTMLFragment(segments []text.Segment, source []byte, literalElement *string) error {
	var fragment []byte
	for _, segment := range segments {
		fragment = append(fragment, source[segment.Start:segment.Stop]...)
	}
	z := htmltoken.NewTokenizer(bytes.NewReader(fragment))
	offset, span, spanStart := 0, 0, 0
	for {
		kind := z.Next()
		if kind == htmltoken.ErrorToken {
			if z.Err() == io.EOF {
				return nil
			}
			return z.Err()
		}
		from := offset
		offset += len(z.Raw())
		if kind != htmltoken.StartTagToken && kind != htmltoken.SelfClosingTagToken && kind != htmltoken.EndTagToken {
			continue
		}
		token := z.Token()
		if *literalElement != "" {
			if kind == htmltoken.EndTagToken && token.Data == *literalElement {
				*literalElement = ""
			}
			continue
		}
		if kind == htmltoken.EndTagToken {
			continue
		}
		switch token.Data {
		case "pre", "code", "script", "style", "textarea", "title":
			*literalElement = token.Data
			continue
		case "img":
		default:
			continue
		}
		for _, attr := range token.Attr {
			if attr.Key != "src" {
				continue
			}
			first, last := htmlImageSource(fragment[from:offset])
			if first < 0 || attr.Val == "" {
				break
			}
			first, last = first+from, last+from
			// Keep the Markdown's list and quote prefixes intact for multiline tags.
			for span < len(segments) && first >= spanStart+segments[span].Stop-segments[span].Start {
				spanStart += segments[span].Stop - segments[span].Start
				span++
			}
			if span == len(segments) || last > spanStart+segments[span].Stop-segments[span].Start {
				break
			}
			url, err := p.resolve(attr.Val)
			if err != nil {
				return err
			}
			if url != attr.Val {
				value := html.EscapeString(url)
				if previous := fragment[first-1]; previous != '"' && previous != '\'' {
					value = strings.NewReplacer("=", "&#61;", "\x60", "&#96;").Replace(value)
				}
				origin := segments[span].Start - spanStart
				p.edits = append(p.edits, imageEdit{start: origin + first, end: origin + last, value: value})
			}
			break // HTML uses the first occurrence of a duplicate attribute.
		}
	}
}

// Consume attributes in order, so src-like text in another attribute is never edited.
var rawHTMLAttr = regexp.MustCompile(`^[\t\n\f\r ]+([^\t\n\f\r />=]+)(?:[\t\n\f\r ]*=[\t\n\f\r ]*(?:"([^"]*)"|'([^']*)'|([^\t\n\f\r >]+)))?`)

func htmlImageSource(raw []byte) (int, int) {
	for offset := len("<img"); offset < len(raw); {
		index := rawHTMLAttr.FindSubmatchIndex(raw[offset:])
		if index == nil {
			break
		}
		if bytes.EqualFold(raw[offset+index[2]:offset+index[3]], []byte("src")) {
			for _, pair := range []int{4, 6, 8} {
				if index[pair] != -1 {
					return offset + index[pair], offset + index[pair+1]
				}
			}
			break
		}
		offset += index[1]
	}
	return -1, -1
}
