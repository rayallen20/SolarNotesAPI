package articlecontent

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRawHTMLAndObservedImages(t *testing.T) {
	root := t.TempDir()
	const base = "http://192.168.1.60:8080"
	reader, err := New(root, base)
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{
		`![查询请求执行过程](./img/查询请求执行过程.jpg)`,
		`![客户端与服务器建立连接后服务器维护的变量](./img/客户端与服务器建立连接后服务器维护的变量.jpg)`,
		`![parentheses](./img/a\(b\).png 'caption "ok"')`,
		`Inline <IMG alt='src="./unrelated.png"' SRC = 'img/a b.png?x=1&amp;y=2' title="icon"> end`,
		"",
		`<div class="gallery">`,
		`  <img`,
		`     src="./img/c.png" data-src="./keep.png" />`,
		`  <img src=img/d.jpg width=80 alt="preview">`,
		`  <img src="https://cdn.example/a.png" srcset="./small.jpg 1x, ./big.jpg 2x">`,
		`  <a href="img/a.png" title="<img src='img/z.png'>">normal</a>`,
		`  <!-- <img src="img/comment.png"> -->`,
		`  <pre><code><img src="img/code.png"></code></pre>`,
		`  <script>const s = '<img src="img/fake.png">'</script>`,
		`</div>`,
		"",
		"> <div>",
		"> <img alt='hello'",
		">   src='img/quote.png'>",
		"> </div>",
		"",
		"`<img src='img/inline-code.png'>`",
		"",
		"```html", `<img src="./img/fenced.png">`, "```",
	}
	want := append([]string{}, lines...)
	for i, filename := range []string{"查询请求执行过程.jpg", "客户端与服务器建立连接后服务器维护的变量.jpg"} {
		want[i] = strings.Replace(lines[i], "./img/"+filename, base+"/img/"+url.PathEscape(filename), 1)
	}
	want[2] = `![parentheses](` + base + `/img/a%28b%29.png 'caption "ok"')`
	want[3] = strings.Replace(lines[3], "img/a b.png", base+"/img/a%20b.png", 1)
	want[7] = strings.Replace(lines[7], "./img/c.png", base+"/img/c.png", 1)
	want[8] = strings.Replace(lines[8], "img/d.jpg", base+"/img/d.jpg", 1)
	want[18] = strings.Replace(lines[18], "img/quote.png", base+"/img/quote.png", 1)
	if err := os.WriteFile(filepath.Join(root, "article.md"), []byte(strings.Join(lines, "\n")), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := reader.ReadMarkdown("article.md")
	if err != nil {
		t.Fatal(err)
	}
	if expected := strings.Join(want, "\n"); got != expected {
		t.Fatalf("unexpected rewritten document:\n%s\nwant:\n%s", got, expected)
	}
}
