package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseArgs(args []string, bi BuildInfo) (Config, error) {
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
		case low == "-m" || low == "--mute":
			cfg.Mute = true
			i++
		case low == "-fm" || low == "--fully-mute":
			cfg.FullyMute = true
			cfg.Mute = true
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
			// first non-option > subject
			goto subject
		}
	}

subject:
	if i >= len(args) {
		return cfg, usageErr("missing required subject (FILE|DIR|images|videos)")
	}

	cfg.SubjectRaw = args[i]
	i++

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
		case low == "-m" || low == "--mute":
			cfg.Mute = true
			i++
		case low == "-fm" || low == "--fully-mute":
			cfg.FullyMute = true
			cfg.Mute = true
			i++
		case strings.HasPrefix(low, "in="):
			if seenTo {
				return cfg, usageErr("in=... must appear before 'to'")
			}
			cfg.InDir = mustAbs(a[3:])
			i++
		case strings.HasPrefix(low, "out="):
			cfg.OutDir = mustAbs(a[4:])
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
	case "images":
		cfg.Subject = SubjectImages
	case "videos":
		cfg.Subject = SubjectVideos
	default:
		fi, err := os.Stat(cfg.SubjectRaw)
		if err != nil {
			return cfg, usageErr("subject not found: " + cfg.SubjectRaw)
		}
		if fi.IsDir() {
			cfg.Subject = SubjectDirImages
			cfg.InDir = mustAbs(cfg.SubjectRaw)
		} else {
			cfg.Subject = SubjectFile
			cfg.FilePath = mustAbs(cfg.SubjectRaw)
			cfg.InDir = mustAbs(filepath.Dir(cfg.FilePath))
		}
	}

	// validate outdir exists
	if st, err := os.Stat(cfg.OutDir); err != nil || !st.IsDir() {
		return cfg, usageErr("output directory not found: " + cfg.OutDir)
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

type usageErr string

func (e usageErr) Error() string { return string(e) }

func usageErrf(f string, a ...any) error { return usageErr(fmt.Sprintf(f, a...)) }
func usageErr(s string) error            { return usageErr(s) }

func mustAbs(p string) string {
	// accept quoted values naturally (shell strips quotes before passing args).
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func IsUsageErr(err error) bool {
	var ue usageErr
	return errors.As(err, &ue)
}
