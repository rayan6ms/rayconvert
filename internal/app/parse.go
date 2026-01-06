package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rayan6ms/rayconvert/internal/mime"
)

func ParseArgs(args []string, bi BuildInfo) (Config, error) {
	outExplicit := false
	cfg := Config{
		InDir:  mustAbs("."),
		OutDir: mustAbs("."),
		Build:  bi,
	}

	if len(args) == 0 {
		return cfg, usageErr("missing arguments")
	}

	// global options may appear before the subject
	i := 0
	for i < len(args) {
		a := args[i]
		low := strings.ToLower(a)

		switch {
		case low == "--help" || low == "-h":
			cfg.Help = true
			return cfg, nil
		case low == "--version":
			cfg.Version = true
			return cfg, nil

		case low == "-ap" || low == "--append":
			cfg.Append = true
			i++
		case low == "-r" || low == "--replace":
			cfg.Replace = true
			i++
		case low == "-m" || low == "--mute":
			cfg.Mute = true
			i++
		case low == "-fm" || low == "--fully-mute":
			cfg.FullyMute = true
			cfg.Mute = true
			i++

		case strings.HasPrefix(low, "-in="):
			cfg.InDir = mustAbs(a[4:])
			i++

		case low == "-in":
			if i+1 >= len(args) {
				return cfg, usageErr("-in requires a value")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "-out="):
			cfg.OutDir = mustAbs(a[5:])
			outExplicit = true
			i++

		case low == "-out":
			if i+1 >= len(args) {
				return cfg, usageErr("-out requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			outExplicit = true
			i += 2
		case low == "-o":
			if i+1 >= len(args) {
				return cfg, usageErr("-o requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			outExplicit = true
			i += 2

		case strings.HasPrefix(low, "-o="):
			cfg.OutDir = mustAbs(a[3:])
			outExplicit = true
			i++

		case low == "-i":
			if i+1 >= len(args) {
				return cfg, usageErr("-i requires a value")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "-i="):
			cfg.InDir = mustAbs(a[3:])
			i++

		case strings.HasPrefix(low, "in="):
			cfg.InDir = mustAbs(a[3:])
			i++
		case strings.HasPrefix(low, "out="):
			cfg.OutDir = mustAbs(a[4:])
			outExplicit = true
			i++

		case strings.HasPrefix(low, "--input="):
			cfg.InDir = mustAbs(a[len("--input="):])
			i++
		case low == "--input":
			if i+1 >= len(args) {
				return cfg, usageErr("--input requires a value")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "--output="):
			cfg.OutDir = mustAbs(a[len("--output="):])
			outExplicit = true
			i++
		case low == "--output":
			if i+1 >= len(args) {
				return cfg, usageErr("--output requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "-"):
			return cfg, usageErr("unknown option: " + a)

		default:
			goto subject
		}
	}

subject:
	if i >= len(args) {
		return cfg, usageErr("missing required subject or 'to'")
	}

	if strings.EqualFold(args[i], "to") {
		cfg.SubjectRaw = ""
	} else {
		cfg.SubjectRaw = args[i]
		i++
	}

	seenTo := false
	for i < len(args) {
		a := args[i]
		low := strings.ToLower(a)

		switch {
		case low == "--help" || low == "-h":
			cfg.Help = true
			return cfg, nil
		case low == "--version":
			cfg.Version = true
			return cfg, nil
		case low == "-ap" || low == "--append":
			cfg.Append = true
			i++
		case low == "-r" || low == "--replace":
			cfg.Replace = true
			i++
		case low == "-m" || low == "--mute":
			cfg.Mute = true
			i++
		case low == "-fm" || low == "--fully-mute":
			cfg.FullyMute = true
			cfg.Mute = true
			i++
		case strings.HasPrefix(low, "-in="):
			if seenTo {
				return cfg, usageErr("-in=... must appear before 'to'")
			}
			cfg.InDir = mustAbs(a[4:])
			i++

		case low == "-in":
			if i+1 >= len(args) {
				return cfg, usageErr("-in requires a value")
			}
			if seenTo {
				return cfg, usageErr("-in must appear before 'to'")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "-out="):
			cfg.OutDir = mustAbs(a[5:])
			outExplicit = true
			i++

		case low == "-out":
			if i+1 >= len(args) {
				return cfg, usageErr("-out requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			outExplicit = true
			i += 2
		case low == "-o":
			if i+1 >= len(args) {
				return cfg, usageErr("-o requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			outExplicit = true
			i += 2

		case strings.HasPrefix(low, "-o="):
			cfg.OutDir = mustAbs(a[3:])
			outExplicit = true
			i++

		case low == "-i":
			if i+1 >= len(args) {
				return cfg, usageErr("-i requires a value")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2

		case strings.HasPrefix(low, "-i="):
			cfg.InDir = mustAbs(a[3:])
			i++

		case strings.HasPrefix(low, "in="):
			if seenTo {
				return cfg, usageErr("in=... must appear before 'to'")
			}
			cfg.InDir = mustAbs(a[3:])
			i++
		case strings.HasPrefix(low, "out="):
			cfg.OutDir = mustAbs(a[4:])
			outExplicit = true
			i++
		case strings.HasPrefix(low, "--input="):
			if seenTo {
				return cfg, usageErr("--input must appear before 'to'")
			}
			cfg.InDir = mustAbs(a[len("--input="):])
			i++
		case low == "--input":
			if i+1 >= len(args) {
				return cfg, usageErr("--input requires a value")
			}
			if seenTo {
				return cfg, usageErr("--input must appear before 'to'")
			}
			cfg.InDir = mustAbs(args[i+1])
			i += 2
		case strings.HasPrefix(low, "--output="):
			cfg.OutDir = mustAbs(a[len("--output="):])
			i++
		case low == "--output":
			if i+1 >= len(args) {
				return cfg, usageErr("--output requires a value")
			}
			cfg.OutDir = mustAbs(args[i+1])
			i += 2
		case low == "to":
			seenTo = true
			i++
		default:
			if !seenTo {
				if strings.EqualFold(a, "images") || strings.EqualFold(a, "videos") {
					return cfg, usageErr("unexpected '" + a + "' before 'to'. Hint: subject must be first (e.g. `rayconvert images in=DIR to jpg`) or use `--input DIR` before subject.")
				}
				return cfg, usageErr("unexpected argument before 'to': " + a)
			}
			if cfg.ToFormat == "" {
				cfg.ToFormat = normalizeFormat(a)
				i++
			} else {
				return cfg, usageErr("unexpected extra argument after format: " + a)
			}
		}
	}

	if !seenTo {
		return cfg, usageErr("missing required keyword: to")
	}
	if cfg.ToFormat == "" {
		return cfg, usageErr("missing required output FORMAT after 'to'")
	}

	// determine subject kind
	subLow := strings.ToLower(cfg.SubjectRaw)
	switch subLow {
	case "":
		if mime.IsVideoFormat(cfg.ToFormat) {
			cfg.Subject = SubjectVideos
		} else if mime.IsImageFormat(cfg.ToFormat) {
			cfg.Subject = SubjectImages
		} else {
			return cfg, usageErr("missing required subject (FILE|DIR|images|videos) and cannot infer from output format: " + cfg.ToFormat)
		}
	case "images":
		cfg.Subject = SubjectImages
	case "videos":
		cfg.Subject = SubjectVideos
	default:
		fi, err := os.Stat(cfg.SubjectRaw)
		if err == nil {
			if fi.IsDir() {
				cfg.Subject = SubjectDirImages
				cfg.InDir = mustAbs(cfg.SubjectRaw)
			} else {
				cfg.Subject = SubjectFile
				cfg.FilePath = mustAbs(cfg.SubjectRaw)
				cfg.InDir = mustAbs(filepath.Dir(cfg.FilePath))
			}
			break
		}

		inFmt := normalizeFormat(cfg.SubjectRaw)
		if mime.IsVideoFormat(inFmt) {
			cfg.Subject = SubjectVideos
			cfg.InFormat = inFmt
			break
		}
		if mime.IsImageFormat(inFmt) {
			cfg.Subject = SubjectImages
			cfg.InFormat = inFmt
			break
		}
		return cfg, usageErr("subject not found: " + cfg.SubjectRaw)
	}

	if !outExplicit {
		cfg.OutDir = cfg.InDir
	}

	if st, err := os.Stat(cfg.OutDir); err != nil || !st.IsDir() {
		return cfg, usageErrf("output directory not found: %s", cfg.OutDir)
	}

	return cfg, nil
}

func normalizeFormat(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, ".")
	s = strings.ToLower(s)
	if s == "jpeg" {
		return "jpg"
	}
	return s
}

type usageError string

func (e usageError) Error() string { return string(e) }

func usageErrf(f string, a ...any) error { return usageError(fmt.Sprintf(f, a...)) }
func usageErr(s string) error            { return usageError(s) }

func IsUsageErr(err error) bool {
	var ue usageError
	return errors.As(err, &ue)
}

func mustAbs(p string) string {
	p = strings.TrimSpace(p)
	p = expandTilde(p)

	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func expandTilde(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return p
		}
		if p == "~" {
			return home
		}
		return filepath.Join(home, p[2:])
	}
	return p
}
