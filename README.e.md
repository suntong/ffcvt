
# {{.Name}}

{{render "license/shields" . "License" "MIT"}}
{{template "badge/godoc" .}}
{{template "badge/goreport" .}}
{{template "badge/travis" .}}

## {{toc 5}}

## {{.Name}} - ffmpeg convert wrapper tool

## Introduction

- The next-generation codec like [High Efficiency Video codec (HEVC), H.265](https://goo.gl/IZrDH2) or [VP9](https://developers.google.com/media/vp9/) can produce videos visually comparable to H.264's result, but in [about half the file size](https://trac.ffmpeg.org/wiki/Encode/H.265).
- Meanwhile the [Opus](https://goo.gl/BPUkTf) [audio codec](https://goo.gl/IZrDH2) is becoming the best thing ever for compressing audio -- A 64K Opus audio stream is comparable to mp3 files of 128K to 256K bandwidth.
- Such fantastic high efficiency audio/video codec/encoding capability has long been available in `ffmpeg`, but fewer people know it or use it, partly because the `ffmpeg` command line is not that simple for every one.
- The `ffcvt` is designed to take the burden from normal Joe -- All you need to do to encode a video is to give one parameter to `ffcvt`, i.e., the path and file name of the video to be encoded, and `ffcvt` will take care of the rest, using the recommended values for both audio/video encoding to properly encode it for you.
- It can't be more simpler than that. However, beneath the simple surface, `ffcvt` is versatile and powerful enough to allow you to touch every corner of audio/video encoding. There is a huge list of environment variables (or command-line parameters) which will allow you tweak the encoding methods and parameters to exactly what you prefer instead.
- Moreover, to encode a directory full of video files, including under its sub-directories, you need just to give `ffcvt` one single parameter, the directory location, and `ffcvt` will go ahead and encode all video files under that directory, including all its sub-directories as well. 

## Quick Usage

There is a quick usage help that comes with `ffcvt`, produced when it is invoked without any parameters:

### $ {{exec "ffcvt" | color "sh"}}


## Environment Variables

For each `ffcvt` command line parameter, there is a environment variable corresponding to it. For example you can use `export FFCVT_FFMPEG=avconv` to use `avconv` instead of `ffmpeg` (Don't, I use it for my [CommandLineArgs](https://github.com/suntong001/lang/blob/master/lang/Go/src/sys/CommandLineArgs.go) to develop/test `ffcvt` without invoking `ffmpeg` each time). 

## Encoding Help

The detailed guide to choose/provide proper parameters to `ffcvt` have been moved to [wiki](https://github.com/suntong/ffcvt/wiki/). For example,

- [HEVC vs VP9](https://github.com/suntong/ffcvt/wiki/KB:-WebM-(VP9)-Encoding#hevc-vs-vp9)
- [HEVC Preset Method Comparison](https://github.com/suntong/ffcvt/wiki/KB:-HEVC-(x265)-Encoding#preset-method-comparison)
- [The HEVC CRF Comparison](https://github.com/suntong/ffcvt/wiki/KB:-HEVC-(x265)-Encoding#the-crf-comparison)
- [Example 1: YouTube Encoding](https://github.com/suntong/ffcvt/wiki/Example:-YouTube-Encoding)
- [Example 2: Talk Encoding](https://github.com/suntong/ffcvt/wiki/Example:-Talk-Encoding)

Please check them out in the [wiki](https://github.com/suntong/ffcvt/wiki/), and for other documents like "Most used ffmpeg options", "How to crop a video", etc.

## Download/Install

### Using `apt`

The `ffcvt` is now officially in Debian repository, so the installation is now as simple as a `apt install`/`apt-get install`:

    apt install ffcvt

### Download binaries

- The latest binary executables are available under  
https://bintray.com/suntong/bin/{{.Name}}/latest  
as the result of the Continuous-Integration process.
- I.e., they are built right from the source code during _every_ git commit _automatically_ by [travis-ci](https://travis-ci.org/).
- Pick & choose the binary executable that suits your OS and its architecture. E.g., for Linux, it would most probably be the `{{.Name}}-linux-amd64` file. If your OS and its architecture is not available in the download list, please let me know and I'll add it.
- You may want to rename it to a shorter name instead, e.g., `{{.Name}}`, after downloading it.


### Debian package

Debian package _repo_ is available at https://dl.bintray.com/suntong/deb.
The _browse-able_ repo view is at https://bintray.com/suntong/deb.

```
echo "deb [trusted=yes] https://dl.bintray.com/suntong/deb all main" | sudo tee /etc/apt/sources.list.d/suntong-debs.list
sudo apt-get update

sudo chmod 644 /etc/apt/sources.list.d/suntong-debs.list
apt-cache policy {{.Name}}

sudo apt-get install -y {{.Name}}
```

### Install Source

If you prefer to compile and install `ffcvt` from source, although a manual process, it's pretty straightforward and simple.

0. Get the source via `git clone` or [`go get`](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies).
0. Do `cd ffcvt`, then issue [`go build`](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies) without any other parameters.
0. Copy the generated executable somewhere in the PATH

That's it, it's ready to roll. 


## Tools Choices

As suggested before, don't use `avconv`, use `ffmpeg` instead (the `avconv` fork was more for political reasons. I personally believe `ffmpeg` is technically superior although might not be politically).

As for video/movie play back, use [mpv](http://mpv.io/). It is a fork of mplayer2 and MPlayer, and is a true *modern* *all-in-one* movie player that can play ANYTHING, and one of the few movie players being actively developed all the time. Download link is in [mpv.io](http://mpv.io/), from which Ubuntu repo I get my Ubuntu `ffmpeg` package as well. If you are unsatisfied with mpv's simple user interface, check out https://wiki.archlinux.org/index.php/Mpv#Front_ends.

## Author(s)

Tong SUN  
![suntong from cpan.org](https://img.shields.io/badge/suntong-%40cpan.org-lightgrey.svg "suntong from cpan.org")

All patches welcome. 
