package articlecontent

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMissingMountIsAStorageError(t *testing.T) {
	reader, err := New(filepath.Join(t.TempDir(), "unmounted"), "http://192.168.1.60:8080")
	if err != nil {
		t.Fatal(err)
	}
	_, err = reader.ReadMarkdown("article.md")
	if !errors.Is(err, fs.ErrNotExist) || errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected the storage root error, got %v", err)
	}
}

func TestMarkdownImageForms(t *testing.T) {
	root := t.TempDir()
	chapter := filepath.Join(root, "books", "chapter")
	if err := os.MkdirAll(chapter, 0755); err != nil {
		t.Fatal(err)
	}
	const base = "http://192.168.1.60:8080"
	reader, err := New(root, base)
	if err != nil {
		t.Fatal(err)
	}
	input := []string{
		`![inline](./img/plan.png "title")`,
		`![legacy](./img/query plan.jpg)`,
		`![escaped](<../shared/%E4%B8%AD p.png?q=1&n=2#mark>)`,
		`![root](/photos/a.png)`,
		`[normal](img/plan.png) ![external](https://cdn.example/a.png) ![cdn](//cdn.example/b.jpg)`,
		`![via ref][p] [link][p]`,
		"",
		`[p]: img/plan.png "original title"`,
		"",
		"`![example](img/example.png)`",
		"",
		"```md", "![sample](img/sample.png)", "```",
		"",
		`<img src="img/html.png">`,
	}
	want := append([]string{}, input...)
	want[0] = `![inline](` + base + `/books/chapter/img/plan.png "title")`
	want[1] = `![legacy](` + base + `/books/chapter/img/query%20plan.jpg)`
	want[2] = `![escaped](<` + base + `/books/shared/%E4%B8%AD%20p.png?q=1&amp;n=2#mark>)`
	want[3] = `![root](` + base + `/photos/a.png)`
	want[5] = `![via ref](` + base + `/books/chapter/img/plan.png "original title") [link][p]`
	want[len(want)-1] = `<img src="` + base + `/books/chapter/img/html.png">`
	path := filepath.Join(chapter, "article.md")
	content := strings.Join(input, "\r\n")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	for _, stored := range []string{path, filepath.Join("books", "chapter", "article.md")} {
		got, err := reader.ReadMarkdown(stored)
		if err != nil {
			t.Fatal(err)
		}
		expected := strings.Join(want, "\r\n")
		if got != expected {
			t.Fatalf("unexpected rewrite:\n%s\nwant:\n%s", got, expected)
		}
	}
	disk, err := os.ReadFile(path)
	if err != nil || string(disk) != content {
		t.Fatal("stored Markdown was modified")
	}
}

func TestUnavailableSourcesAndImagePaths(t *testing.T) {
	root := t.TempDir()
	reader, err := New(root, "http://192.168.1.60:8080")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bad.md"), []byte{0xff}, 0600); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "private.md"), []byte("other location"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "outside")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "directory.md"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{
		"absent.md", "bad.md", "a.pdf", "../outside.md", "directory.md",
		filepath.Join(outside, "private.md"), filepath.Join("outside", "private.md"),
	} {
		if _, err := reader.ReadMarkdown(file); !errors.Is(err, ErrUnavailable) {
			t.Errorf("read %q: want unavailable, got %v", file, err)
		}
	}
	for _, image := range []string{"../image.png", ".hidden/image.png", "private.MD", "outside/a.png"} {
		for _, content := range []string{"![x](" + image + ")", "<img src='" + image + "'>"} {
			if err := os.WriteFile(filepath.Join(root, "article.md"), []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := reader.ReadMarkdown("article.md"); !errors.Is(err, ErrUnavailable) {
				t.Errorf("image %q: want unavailable, got %v", content, err)
			}
		}
	}
}

func TestSharedReaderResolvesEachArticleDirectory(t *testing.T) {
	root := t.TempDir()
	const base = "http://192.168.1.60:8080"
	reader, err := New(root, base)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"chapter-one", "chapter-two"} {
		t.Run(dir, func(t *testing.T) {
			t.Parallel()
			if err := os.Mkdir(filepath.Join(root, dir), 0755); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "article.md")
			if err := os.WriteFile(filepath.Join(root, path), []byte("![figure](./img/1.jpg)"), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := reader.ReadMarkdown(path)
			if err != nil || got != "![figure]("+base+"/"+dir+"/img/1.jpg)" {
				t.Fatalf("unexpected read for %s: %q, %v", dir, got, err)
			}
		})
	}
}
