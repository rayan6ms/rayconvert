package ff

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rayan6ms/rayconvert/internal/mime"
)

type Converter struct {
	mute      bool
	fullyMute bool
	append    bool
}

func NewConverter(mute, fullyMute, append bool) *Converter {
	if fullyMute {
		mute = true
	}
	return &Converter{mute: mute, fullyMute: fullyMute, append: append}
}

func (c *Converter) ConvertOne(inPath, outDir, outFmt string) (bool, bool, error) {
	inAbs, _ := filepath.Abs(inPath)
	outDirAbs, _ := filepath.Abs(outDir)

	base := strings.TrimSuffix(filepath.Base(inAbs), filepath.Ext(inAbs))
	outFmt = strings.TrimPrefix(strings.TrimSpace(outFmt), ".")
	outPath := filepath.Join(outDirAbs, base+"."+outFmt)

	outAbs, _ := filepath.Abs(outPath)
	if outAbs == inAbs {
		c.msgErr(fmt.Sprintf("skipping '%s' (would overwrite itself)", filepath.Base(inAbs)))
		return false, true, nil
	}

	if strings.TrimPrefix(strings.ToLower(filepath.Ext(inAbs)), ".") == strings.ToLower(outFmt) &&
		filepath.Dir(inAbs) == outDirAbs {
		return false, true, nil
	}

	if c.append {
		outPath = pickUnique(outPath)
	}

	if strings.EqualFold(outFmt, "svg") {
		ok, err := c.convertToSVG(inAbs, outPath)
		if err != nil {
			return false, false, err
		}
		if ok && c.mute && !c.fullyMute {
			fmt.Printf("successfully converted '%s' to '%s'\n", filepath.Base(inAbs), outFmt)
		}
		return ok, false, nil
	}

	isVideoToImage := false
	k := mime.DetectKind(filepath.Dir(inAbs), inAbs)
	if k == mime.KindVideo && mime.IsImageFormat(outFmt) {
		isVideoToImage = true
	}

	args := []string{}
	if c.append {
		args = append(args, "-n")
	} else {
		args = append(args, "-y")
	}

	args = append(args, "-hide_banner")

	args = append(args, "-i", inAbs)

	if isVideoToImage {
		args = append(args, "-frames:v", "1")
	}

	if strings.ToLower(outFmt) == "jpg" {
		args = append(args, "-q:v", "2")
	}

	if mime.IsImageFormat(outFmt) && strings.ToLower(outFmt) != "gif" {
		args = append(args, "-frames:v", "1", "-update", "1")
	}

	args = append(args, outPath)

	if err := c.runTool("ffmpeg", args); err != nil {
		return false, false, err
	}

	if c.mute && !c.fullyMute {
		fmt.Printf("successfully converted '%s' to '%s'\n", filepath.Base(inAbs), outFmt)
	}
	return true, false, nil
}

func (c *Converter) convertToSVG(inAbs, outSVG string) (bool, error) {
	tmpDir, err := os.MkdirTemp("", "rayconvert-svg-*")
	if err != nil {
		c.msgErr("error: failed creating temp dir")
		return false, errors.New("svg conversion failed")
	}
	defer os.RemoveAll(tmpDir)

	tmpBmp := filepath.Join(tmpDir, "input.bmp")

	ffArgs := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-y",
		"-i", inAbs,
		"-frames:v", "1",
		"-update", "1",
		tmpBmp,
	}
	if err := c.runTool("ffmpeg", ffArgs); err != nil {
		return false, err
	}

	ptArgs := []string{"-b", "svg", "-o", outSVG, tmpBmp}
	if err := c.runTool("potrace", ptArgs); err != nil {
		c.msgErr("error: potrace is required for 'to svg' (install it on the host, e.g. `sudo apt install potrace`)")
		return false, errors.New("svg conversion failed")
	}

	return true, nil
}

func (c *Converter) runTool(tool string, args []string) error {
	cmd := toolCmd(tool, args...)

	if c.fullyMute {
		cmd.Stdout = nil
		cmd.Stderr = nil
		return cmd.Run()
	}

	if c.mute {
		out, err := cmd.CombinedOutput()
		if err != nil {
			msg := bestErrorLine(string(out))
			if msg == "" {
				msg = "error running " + tool
			}
			c.msgErr("error: " + msg)
			return errors.New(tool + " failed")
		}
		return nil
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func toolCmd(prog string, args ...string) *exec.Cmd {
	if os.Getenv("FLATPAK_ID") != "" {
		if _, err := exec.LookPath(prog); err == nil {
			return exec.Command(prog, args...)
		}
		return exec.Command("flatpak-spawn", append([]string{"--host", prog}, args...)...)
	}
	return exec.Command(prog, args...)
}

func (c *Converter) msgErr(s string) {
	if !c.fullyMute {
		fmt.Fprintln(os.Stderr, "rayconvert:", s)
	}
}

func bestErrorLine(out string) string {
	lines := strings.Split(out, "\n")
	clean := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		if strings.HasPrefix(ln, "ffmpeg version") ||
			strings.HasPrefix(ln, "configuration:") ||
			strings.HasPrefix(ln, "libav") {
			continue
		}
		clean = append(clean, ln)
	}

	best := ""
	bestScore := -1
	for _, ln := range clean {
		s := scoreLine(ln)
		if s > bestScore {
			bestScore = s
			best = ln
		}
	}

	if bestScore <= 0 && len(clean) > 0 {
		return clean[len(clean)-1]
	}
	return best
}

func scoreLine(ln string) int {
	l := strings.ToLower(ln)
	switch {
	case strings.Contains(l, "unable to choose an output format"):
		return 100
	case strings.Contains(l, "unknown encoder"), strings.Contains(l, "unknown format"):
		return 95
	case strings.Contains(l, "error initializing the muxer"):
		return 90
	case strings.Contains(l, "error opening output file"):
		return 85
	case strings.Contains(l, "no such file or directory"), strings.Contains(l, "permission denied"):
		return 80
	case strings.Contains(l, "invalid argument"):
		return 70
	case strings.Contains(l, "not found"):
		return 60
	case strings.Contains(l, "error"):
		return 50
	default:
		return 1
	}
}

func pickUnique(p string) string {
	if _, err := os.Stat(p); err != nil {
		return p
	}
	dir := filepath.Dir(p)
	base := strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
	ext := filepath.Ext(p)
	for n := 1; ; n++ {
		try := filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, n, ext))
		if _, err := os.Stat(try); err != nil {
			return try
		}
	}
}
