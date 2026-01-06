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

	err := c.runFFmpeg(args)
	if err != nil {
		return false, false, err
	}

	if c.mute && !c.fullyMute {
		fmt.Printf("successfully converted '%s' to '%s'\n", filepath.Base(inAbs), outFmt)
	}
	return true, false, nil
}

func (c *Converter) runFFmpeg(args []string) error {
	cmd := exec.Command("ffmpeg", args...)
	if os.Getenv("FLATPAK_ID") != "" {
		cmd = exec.Command("flatpak-spawn", append([]string{"--host", "ffmpeg"}, args...)...)
	}

	if c.fullyMute {
		cmd.Stdout = nil
		cmd.Stderr = nil
		return cmd.Run()
	}

	if c.mute {
		out, err := cmd.CombinedOutput()
		if err != nil {
			last := lastNonEmptyLine(string(out))
			if last != "" {
				c.msgErr("error: " + last)
			} else {
				c.msgErr("error running ffmpeg")
			}
			return errors.New("ffmpeg failed")
		}
		return nil
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Converter) msgErr(s string) {
	if !c.fullyMute {
		fmt.Fprintln(os.Stderr, "rayconvert:", s)
	}
}

func lastNonEmptyLine(s string) string {
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			return strings.TrimSpace(lines[i])
		}
	}
	return ""
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
