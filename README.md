# rayconvert

A small CLI wrapper around **ffmpeg** that converts a single file or batches of files with a readable syntax and safer defaults.

## Usage

```bash
rayconvert (FILE|DIR|images|videos) [in=DIR] to FORMAT [out=DIR] [-ap|--append] [-m|--mute] [-fm|--fully-mute]
```

- Examples:

    ```bash
      rayconvert . to jpg
    ```

    ```bash
      rayconvert images to png --append
    ```

    ```bash
      rayconvert videos in="./clips" to mp4 out="./out" --mute
    ```

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
