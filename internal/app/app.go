package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rayan6ms/rayconvert/internal/ff"
	"github.com/rayan6ms/rayconvert/internal/mime"
)

func Run(args []string, bi BuildInfo) int {
	cfg, err := ParseArgs(args, bi)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rayconvert:", err)
		printUsage(os.Stderr)
		return 2
	}

	if cfg.Help {
		printUsage(os.Stdout)
		return 0
	}
	if cfg.Version {
		fmt.Printf("rayconvert %s (%s) %s\n", cfg.Build.Version, cfg.Build.Commit, cfg.Build.Date)
		return 0
	}

	// enforce format compatibility for group modes
	if cfg.Subject == SubjectVideos {
		if mime.IsImageFormat(cfg.ToFormat) {
			fmt.Fprintln(os.Stderr, "rayconvert: invalid target: videos to", cfg.ToFormat, "(image format)")
			return 2
		}
		if !mime.IsVideoFormat(cfg.ToFormat) {
			fmt.Fprintln(os.Stderr, "rayconvert: invalid target: videos to", cfg.ToFormat, "(unknown video format)")
			return 2
		}
	} else if cfg.Subject == SubjectImages || cfg.Subject == SubjectDirImages {
		if mime.IsVideoFormat(cfg.ToFormat) {
			fmt.Fprintln(os.Stderr, "rayconvert: invalid target: images to", cfg.ToFormat, "(video format)")
			return 2
		}
		if !mime.IsImageFormat(cfg.ToFormat) {
			fmt.Fprintln(os.Stderr, "rayconvert: invalid target: images to", cfg.ToFormat, "(unknown image format)")
			return 2
		}
	}

	conv := ff.NewConverter(cfg.Mute, cfg.FullyMute, cfg.Append)

	converted, skipped, failed := 0, 0, 0

	convertOne := func(path string) {
		ok, sk, err := conv.ConvertOne(path, cfg.OutDir, cfg.ToFormat)
		if sk {
			skipped++
			return
		}
		if err != nil {
			failed++
			return
		}
		if ok {
			converted++
			if !cfg.Append {
				_ = os.Remove(path)
			}
		}
	}

	switch cfg.Subject {
	case SubjectFile:
		// for single file validate based on detected type
		kind := mime.DetectKind(pathBaseForDetect(cfg.FilePath), cfg.FilePath)
		if kind == mime.KindImage && !mime.IsImageFormat(cfg.ToFormat) {
			fmt.Fprintln(os.Stderr, "rayconvert: file is an image; target format must be an image format")
			failed++
		} else if kind == mime.KindVideo && !(mime.IsVideoFormat(cfg.ToFormat) || mime.IsImageFormat(cfg.ToFormat)) {
			fmt.Fprintln(os.Stderr, "rayconvert: file is a video; target must be a video format or image format (thumbnail)")
			failed++
		} else {
			convertOne(cfg.FilePath)
		}

	default:
		entries, err := os.ReadDir(cfg.InDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "rayconvert:", err)
			return 1
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			p := filepath.Join(cfg.InDir, e.Name())

			if cfg.Subject == SubjectVideos {
				if mime.DetectKind(pathBaseForDetect(p), p) != mime.KindVideo {
					skipped++
					continue
				}
			} else {
				if mime.DetectKind(pathBaseForDetect(p), p) != mime.KindImage {
					skipped++
					continue
				}
			}
			convertOne(p)
		}
	}

	if !cfg.FullyMute && !cfg.Mute {
		fmt.Printf("rayconvert: converted=%d skipped=%d failed=%d (in=%s out=%s append=%v to=%s)\n",
			converted, skipped, failed, cfg.InDir, cfg.OutDir, cfg.Append, cfg.ToFormat)
	}

	if failed > 0 {
		return 1
	}
	return 0
}

func pathBaseForDetect(p string) string { return filepath.Dir(p) }

func printUsage(f *os.File) {
	fmt.Fprintln(f, `Usage:
  rayconvert (FILE|DIR|images|videos) [in=DIR] to FORMAT [out=DIR] [-ap|--append] [-m|--mute] [-fm|--fully-mute]
  rayconvert --help
  rayconvert --version

Notes:
  - "jpeg" is treated as "jpg"
  - DIR is treated as "images in that directory" (so: rayconvert . to jpg)
  - "videos ... to jpg" is rejected (videos require video output formats)
  - in= and out= accept quoted values via your shell (e.g. out="path with spaces")
`)
}
