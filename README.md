# rayconvert

A small CLI wrapper around **ffmpeg** that converts a single file or batches of files with a readable syntax.

## Usage

```bash
rayconvert (FILE|DIR|images|videos) [in=DIR|-i DIR|--input DIR] to FORMAT [out=DIR|-o DIR|--output DIR] [-ap|--append] [-r|--replace] [-m|--mute] [-fm|--fully-mute]
```

## Defaults

- Keeps originals by default (no deletion).

- If out=/--output is not provided, output defaults to the input directory.

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

Grab the latest release from GitHub Releases and put rayconvert somewhere in your PATH.

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

MIT (see LICENSE)
