
# ffcvt

[![MIT License](http://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/suntong/ffcvt?status.svg)](http://godoc.org/github.com/suntong/ffcvt)
[![Go Report Card](https://goreportcard.com/badge/github.com/suntong/ffcvt)](https://goreportcard.com/report/github.com/suntong/ffcvt)
[![travis Status](https://travis-ci.org/suntong/ffcvt.svg?branch=master)](https://travis-ci.org/suntong/ffcvt)

## TOC
- [ffcvt - ffmpeg convert wrapper tool](#ffcvt---ffmpeg-convert-wrapper-tool)
- [Introduction](#introduction)
- [Quick Usage](#quick-usage)
  - [$ ffcvt](#-ffcvt)
- [Environment Variables](#environment-variables)
- [Encoding Help](#encoding-help)
- [Download/Install](#downloadinstall)
  - [Using `apt`](#using-`apt`)
  - [Download binaries](#download-binaries)
  - [Debian package](#debian-package)
  - [Install Source](#install-source)
- [Tools Choices](#tools-choices)
- [Author(s)](#author(s))

## ffcvt - ffmpeg convert wrapper tool

## Introduction

- The next-generation codec like [High Efficiency Video codec (HEVC), H.265](https://goo.gl/IZrDH2) or [VP9](https://developers.google.com/media/vp9/) can produce videos visually comparable to H.264's result, but in [about half the file size](https://trac.ffmpeg.org/wiki/Encode/H.265).
- Meanwhile the [Opus](https://goo.gl/BPUkTf) [audio codec](https://goo.gl/IZrDH2) is becoming the best thing ever for compressing audio -- A 64K Opus audio stream is comparable to mp3 files of 128K to 256K bandwidth.
- Such fantastic high efficiency audio/video codec/encoding capability has long been available in `ffmpeg`, but fewer people know it or use it, partly because the `ffmpeg` command line is not that simple for every one.
- The `ffcvt` is designed to take the burden from normal Joe -- All you need to do to encode a video is to give one parameter to `ffcvt`, i.e., the path and file name of the video to be encoded, and `ffcvt` will take care of the rest, using the recommended values for both audio/video encoding to properly encode it for you.
- It can't be more simpler than that. However, beneath the simple surface, `ffcvt` is versatile and powerful enough to allow you to touch every corner of audio/video encoding. There is a huge list of environment variables (or command-line parameters) which will allow you tweak the encoding methods and parameters to exactly what you prefer instead.
- Moreover, to encode a directory full of video files, including under its sub-directories, you need just to give `ffcvt` one single parameter, the directory location, and `ffcvt` will go ahead and encode all video files under that directory, including all its sub-directories as well. 

## Quick Usage

There is a quick usage help that comes with `ffcvt`, produced when it is invoked without any parameters:

### $ ffcvt
```sh
Usage:
 ffcvt [flags] 

Flags:

  -t	target type: webm/x265-opus/x264-mp3/wx/youtube (FFCVT_T)
  -ves	video encoding method set (FFCVT_VES)
  -aes	audio encoding method set (FFCVT_AES)
  -ses	subtitle encoding method set (FFCVT_SES)
  -vep	video encoding method prepend (FFCVT_VEP)
  -aep	audio encoding method prepend (FFCVT_AEP)
  -sep	subtitle encoding method prepend (FFCVT_SEP)
  -vea	video encoding method append (FFCVT_VEA)
  -aea	audio encoding method append (FFCVT_AEA)
  -abr	audio bitrate (64k for opus, 256k for mp3) (FFCVT_ABR)
  -crf	the CRF value: 0-51. Higher CRF gives lower quality
	 (28 for x265, ~ 23 for x264) (FFCVT_CRF)

  -d	directory that hold input files (FFCVT_D)
  -f	input file name (either -d or -f must be specified) (FFCVT_F)
  -sym	symlinks will be processed as well (FFCVT_SYM)
  -exts	extension list for all the files to be queued (FFCVT_EXTS)
  -suf	suffix to the output file names (FFCVT_SUF)
  -ext	extension for the output file (FFCVT_EXT)
  -w	work directory that hold output files (FFCVT_W)

  -ac	copy audio codec (FFCVT_AC)
  -vc	copy video codec (FFCVT_VC)
  -an	no audio, output video only (FFCVT_AN)
  -vn	no video, output audio only (FFCVT_VN)
  -vss	video: same size (FFCVT_VSS)
  -lang	language selection for audio stream extraction (FFCVT_LANG)
  -o	more options that will pass to ffmpeg program (FFCVT_O)
  -ato-opus	audio encode to opus, using -abr (FFCVT_ATO_OPUS)
  -vto-x265	video video encode to x265, using -crf (FFCVT_VTO_X265)

  -p	par2create, create par2 files (in work directory) (FFCVT_P)
  -nc	no clobber, do not queue those already been converted (FFCVT_NC)
  -n	no exec, dry run (FFCVT_N)

  -force	overwrite any existing none-empty file (FFCVT_FORCE)
  -debug	debugging level (FFCVT_DEBUG)
  -ffmpeg	ffmpeg program executable name (FFCVT_FFMPEG)
  -ffprobe	ffprobe program execution (FFCVT_FFPROBE)
  -version	print version then exit (FFCVT_VERSION)

Details:

  -abr string
    	audio bitrate (64k for opus, 256k for mp3)
  -ac
    	copy audio codec
  -aea string
    	audio encoding method append
  -aep string
    	audio encoding method prepend
  -aes string
    	audio encoding method set
  -an
    	no audio, output video only
  -ato-opus
    	audio encode to opus, using -abr
  -crf string
    	the CRF value: 0-51. Higher CRF gives lower quality
    		 (28 for x265, ~ 23 for x264)
  -d string
    	directory that hold input files
  -debug int
    	debugging level (default 1)
  -ext string
    	extension for the output file
  -exts string
    	extension list for all the files to be queued (default ".3GP.3G2.ASF.AVI.DAT.DIVX.FLV.M2TS.M4V.MKV.MOV.MPEG.MP4.MPG.RMVB.RM.TS.VOB.WEBM.WMV")
  -f string
    	input file name (either -d or -f must be specified)
  -ffmpeg string
    	ffmpeg program executable name (default "ffmpeg")
  -ffprobe string
    	ffprobe program execution (default "ffprobe -print_format flat")
  -force
    	overwrite any existing none-empty file
  -lang string
    	language selection for audio stream extraction (default "eng")
  -n	no exec, dry run
  -nc
    	no clobber, do not queue those already been converted
  -o string
    	more options that will pass to ffmpeg program
  -p	par2create, create par2 files (in work directory)
  -sep string
    	subtitle encoding method prepend
  -ses string
    	subtitle encoding method set
  -suf string
    	suffix to the output file names
  -sym
    	symlinks will be processed as well
  -t string
    	target type: webm/x265-opus/x264-mp3/wx/youtube (default "webm")
  -vc
    	copy video codec
  -vea string
    	video encoding method append
  -vep string
    	video encoding method prepend
  -version
    	print version then exit
  -ves string
    	video encoding method set
  -vn
    	no video, output audio only
  -vss
    	video: same size (default true)
  -vto-x265
    	video video encode to x265, using -crf
  -w string
    	work directory that hold output files

To reduce output, use `-debug 0`, e.g., `ffcvt -force -debug 0 -f testf.mp4 ...`
```


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
https://bintray.com/suntong/bin/ffcvt/latest  
as the result of the Continuous-Integration process.
- I.e., they are built right from the source code during _every_ git commit _automatically_ by [travis-ci](https://travis-ci.org/).
- Pick & choose the binary executable that suits your OS and its architecture. E.g., for Linux, it would most probably be the `ffcvt-linux-amd64` file. If your OS and its architecture is not available in the download list, please let me know and I'll add it.
- You may want to rename it to a shorter name instead, e.g., `ffcvt`, after downloading it.


### Debian package

Debian package _repo_ is available at https://dl.bintray.com/suntong/deb.
The _browse-able_ repo view is at https://bintray.com/suntong/deb.

```
echo "deb [trusted=yes] https://dl.bintray.com/suntong/deb all main" | sudo tee /etc/apt/sources.list.d/suntong-debs.list
sudo apt-get update

sudo chmod 644 /etc/apt/sources.list.d/suntong-debs.list
apt-cache policy ffcvt

sudo apt-get install -y ffcvt
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
