package mime

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Kind int

const (
	KindUnknown Kind = iota
	KindImage
	KindVideo
)

func IsImageFormat(fmt string) bool {
	fmt = strings.ToLower(fmt)
	switch fmt {
	case "jpg", "jpeg", "png", "webp", "bmp", "tif", "tiff", "gif", "avif", "heic", "heif":
		return true
	default:
		return false
	}
}

func IsVideoFormat(fmt string) bool {
	fmt = strings.ToLower(fmt)
	switch fmt {
	case "mp4", "mov", "mkv", "webm", "avi", "m4v", "mpg", "mpeg", "ts", "mts", "m2ts", "3gp", "ogv":
		return true
	default:
		return false
	}
}

func DetectKind(_baseDir, path string) Kind {
	if mt, ok := mimeTypeViaFile(path); ok {
		if strings.HasPrefix(mt, "image/") {
			return KindImage
		}
		if strings.HasPrefix(mt, "video/") {
			return KindVideo
		}
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if IsImageFormat(ext) {
		return KindImage
	}
	if IsVideoFormat(ext) {
		return KindVideo
	}
	return KindUnknown
}

func mimeTypeViaFile(path string) (string, bool) {
	if os.Getenv("FLATPAK_ID") != "" {
		out, err := exec.Command("flatpak-spawn", "--host", "file", "-b", "--mime-type", "--", path).CombinedOutput()
		if err == nil {
			return strings.TrimSpace(string(out)), true
		}
	}

	out, err := exec.Command("file", "-b", "--mime-type", "--", path).CombinedOutput()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(bytes.TrimSpace(out))), true
}
