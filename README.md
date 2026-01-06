# rayconvert

A small CLI wrapper around **ffmpeg** that converts a single file or batches of files with a readable syntax.

You can also read the manual with:

```bash
man rayconvert
```

## Usage for ~~dorks~~ *Flatpak* users

If you installed rayconvert via Flatpak, you normally run it like this:

```bash
flatpak run io.github.rayan6ms.Rayconvert ...
```

To use the simpler `rayconvert` command, you can either create an alias in your shell configuration file (e.g. `~/.bashrc` or `~/.zshrc`):

```bash
alias rayconvert='flatpak run io.github.rayan6ms.Rayconvert'
```

Or you can create a small wrapper script (recommended) in `~/.local/bin`:

```bash
install -d ~/.local/bin
cat > ~/.local/bin/rayconvert <<'EOF'
#!/usr/bin/env bash
exec flatpak run io.github.rayan6ms.Rayconvert "$@"
EOF
chmod +x ~/.local/bin/rayconvert
```

Make sure `~/.local/bin` is in your PATH (then restart your shell):

```bash
echo "$PATH" | grep -q "$HOME/.local/bin" || echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

(If you use zsh, replace `~/.bashrc` with `~/.zshrc`.)

If you prefer a system-wide wrapper:

```bash
sudo install -d /usr/local/bin
sudo tee /usr/local/bin/rayconvert >/dev/null <<'EOF'
#!/usr/bin/env bash
exec flatpak run io.github.rayan6ms.Rayconvert "$@"
EOF
sudo chmod +x /usr/local/bin/rayconvert
```

## Usage

```bash
rayconvert (FILE|DIR|images|videos) [INPUT] to FORMAT [OUTPUT] [FLAGS]
```

### Subject (required)

The first non-option argument must be one of:

* `FILE` — convert a single file
* `DIR` — treated as `images` for that directory (so `rayconvert . to jpg` is valid)
* `images` — convert all images in the input directory
* `videos` — convert all videos in the input directory

### INPUT (optional)

All of these are supported (case-insensitive for `in=`):

* `in=DIR`
* `-i DIR`
* `-i=DIR`
* `--input DIR`
* `--input=DIR`
* `-in DIR`
* `-in=DIR`

Paths can be quoted, and `~` is expanded:

```bash
in="~/Pictures"
-in="~/Pictures"
--input="~/Pictures"
```

### OUTPUT (optional)

All of these are supported (case-insensitive for `out=`):

* `out=DIR`
* `-o DIR`
* `-o=DIR`
* `--output DIR`
* `--output=DIR`
* `-out DIR`
* `-out=DIR`

If output is not provided, it defaults to the **input directory**.

### FLAGS (optional)

* `-ap`, `--append`
  If an output filename already exists, choose a unique name instead of overwriting it (e.g. `name-1.jpg`).

* `-r`, `--replace`
  Delete the original file after a successful conversion. (Default is **keep originals**.)

* `-m`, `--mute`
  Hide ffmpeg output and print simplified messages.

* `-fm`, `--fully-mute`
  Print nothing at all (implies `--mute`). Exit code still indicates success/failure.

## Defaults

* Keeps originals by default (no deletion). Use `-r` / `--replace` to delete originals after success.
* If output is not provided, output defaults to the input directory.

## Examples

Convert images in the current directory to jpg (keeps originals):

```bash
rayconvert . to jpg
```

Convert all images in a directory (keeps originals), simplified output:

```bash
rayconvert images in="$HOME/Pictures/quizzes/test" to jpg -m
```

Replace originals (delete after successful conversion):

```bash
rayconvert images in="$HOME/Pictures/quizzes/test" to jpg -m --replace
```

Convert videos to mp4 into a different output folder:

```bash
rayconvert videos in="./clips" to mp4 out="./out" --mute
```

Convert a single file:

```bash
rayconvert "./My File.mov" to jpg --append
```

## Install

### Option A: Download a binary (recommended)

Grab the latest release from GitHub Releases and put `rayconvert` somewhere in your PATH.

### Option B: Build from source

Requires Go and ffmpeg.

```bash
git clone https://github.com/rayan6ms/rayconvert
cd rayconvert
go build -o rayconvert ./cmd/rayconvert
sudo install -m755 rayconvert /usr/local/bin/rayconvert
sudo install -m644 man/rayconvert.1 /usr/local/share/man/man1/rayconvert.1
```

## License

MIT (see [LICENSE](LICENSE))

