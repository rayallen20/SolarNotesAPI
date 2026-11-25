package articlecontent

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf8"
)

// ErrUnavailable identifies unavailable Markdown files and invalid local image references.
var ErrUnavailable = errors.New("article content unavailable")

// Reader shares immutable settings across requests.
type Reader struct {
	rootDir string
	base    *url.URL
}

func New(rootDir, publicBaseURL string) (*Reader, error) {
	if !filepath.IsAbs(rootDir) {
		return nil, errors.New("storage.root_dir must be absolute")
	}
	base, err := url.Parse(publicBaseURL)
	if err != nil || (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" ||
		base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" {
		return nil, errors.New("asset.base_url must be an absolute HTTP(S) resource base")
	}
	return &Reader{rootDir: filepath.Clean(rootDir), base: base}, nil
}

type articleSource struct {
	root       *os.Root
	assetBase  *url.URL
	articleDir string
}

// ReadMarkdown opens a fresh storage root for this request and keeps the original Markdown layout.
func (r *Reader) ReadMarkdown(storedPath string) (string, error) {
	realRoot, err := filepath.EvalSymlinks(r.rootDir)
	if err != nil {
		return "", err
	}
	root, err := os.OpenRoot(realRoot)
	if err != nil {
		return "", err
	}
	defer root.Close()

	source := articleSource{root: root, assetBase: r.base}
	data, err := source.readMarkdown(r.rootDir, storedPath)
	if err != nil {
		return "", err
	}
	return RewriteImages(data, source.imageURL)
}

// readMarkdown 读取 UTF-8 正文, 成功时记下文章所在目录相对于存储根的真实路径.
func (a *articleSource) readMarkdown(rootDir, storedPath string) ([]byte, error) {
	name := storedPath
	if filepath.IsAbs(name) {
		relative, err := filepath.Rel(rootDir, name)
		if err != nil {
			return nil, ErrUnavailable
		}
		name = relative
	}
	if !strings.EqualFold(filepath.Ext(name), ".md") {
		return nil, ErrUnavailable
	}
	actual, err := a.resolveStoragePath(name, false)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(filepath.Ext(actual), ".md") {
		return nil, ErrUnavailable
	}
	// Linux 和 macOS 上用 O_NONBLOCK 打开, 避免在 Stat 识别 FIFO 之前阻塞.
	f, err := a.root.OpenFile(actual, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrUnavailable
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, ErrUnavailable
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(content) {
		return nil, ErrUnavailable
	}
	a.articleDir = filepath.Dir(actual)
	return content, nil
}

// resolveStoragePath 解析符号链接并检查根边界. 普通缺图可以保留尚不存在的尾部路径.
func (a *articleSource) resolveStoragePath(name string, allowMissing bool) (string, error) {
	if !filepath.IsLocal(name) || strings.ContainsRune(name, 0) {
		return "", ErrUnavailable
	}
	rootDir := a.root.Name()
	candidate := filepath.Join(rootDir, name)
	var missing []string
	for {
		resolved, err := filepath.EvalSymlinks(candidate)
		if err == nil {
			relative, relErr := filepath.Rel(rootDir, resolved)
			if relErr != nil || !filepath.IsLocal(relative) {
				return "", ErrUnavailable
			}
			for i := len(missing) - 1; i >= 0; i-- {
				relative = filepath.Join(relative, missing[i])
			}
			return relative, nil
		}
		if errors.Is(err, syscall.ENOTDIR) || errors.Is(err, syscall.ELOOP) {
			return "", ErrUnavailable
		}
		if !errors.Is(err, fs.ErrNotExist) || candidate == rootDir {
			return "", err
		}
		if !allowMissing {
			return "", ErrUnavailable
		}
		// 失效的符号链接本身仍有目录项, 普通缺失路径则没有.
		_, statErr := os.Lstat(candidate)
		if statErr == nil {
			return "", ErrUnavailable
		}
		if !errors.Is(statErr, fs.ErrNotExist) {
			return "", statErr
		}
		missing = append(missing, filepath.Base(candidate))
		candidate = filepath.Dir(candidate)
	}
}

// imageURL 将文章引用的本地图片映射到该文件的公开 URL.
func (a *articleSource) imageURL(destination string) (string, error) {
	u, err := url.Parse(destination)
	if err != nil {
		return "", fmt.Errorf("%w: invalid image URI", ErrUnavailable)
	}
	if u.Host != "" && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "") {
		return destination, nil
	}
	if u.Scheme != "" || u.Path == "" || strings.Contains(u.Path, "\\") {
		return "", ErrUnavailable
	}
	name := filepath.FromSlash(u.Path)
	if strings.HasPrefix(u.Path, "/") {
		name = strings.TrimPrefix(name, string(filepath.Separator))
	} else {
		name = filepath.Join(a.articleDir, name)
	}
	if !a.isPublicImagePath(filepath.Clean(name)) {
		return "", ErrUnavailable
	}
	actual, err := a.resolveStoragePath(name, true)
	if err != nil {
		return "", err
	}
	if !a.isPublicImagePath(actual) {
		return "", ErrUnavailable
	}
	path := filepath.ToSlash(actual)
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	result := *a.assetBase
	result.Path = strings.TrimRight(result.Path, "/") + "/" + path
	result.RawPath = strings.TrimRight(a.assetBase.EscapedPath(), "/") + "/" + strings.Join(segments, "/")
	result.RawQuery, result.ForceQuery = u.RawQuery, u.ForceQuery
	result.Fragment, result.RawFragment = u.Fragment, u.RawFragment
	return result.String(), nil
}

// isPublicImagePath 检查 Nginx 对隐藏路径分量和 .md 的限制.
func (a *articleSource) isPublicImagePath(name string) bool {
	if !filepath.IsLocal(name) || strings.EqualFold(filepath.Ext(name), ".md") {
		return false
	}
	for _, segment := range strings.Split(filepath.ToSlash(name), "/") {
		if strings.HasPrefix(segment, ".") {
			return false
		}
	}
	return true
}
