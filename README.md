# ffcvt - ffmpeg convert wrapper tool

## Introduction

- The next-generation [High Efficiency Video codec (HEVC), H.265](https://goo.gl/IZrDH2) can produce videos visually comparable to libx264's result, but in [about half the file size](https://trac.ffmpeg.org/wiki/Encode/H.265).
- Meanwhile the [Opus](https://goo.gl/BPUkTf) [audio codec](https://goo.gl/IZrDH2) is becoming the best thing ever for compressing audio -- A 64K Opus audio stream is comparable to mp3 files of 128K to 256K bandwidth.
- Such fantastic high efficiency audio/video codec/encoding capability has long been available in `ffmpeg`, but fewer people know it or use it, partly because the `ffmpeg` command line is not that simple for every one.
- The `ffcvt` is designed to take the burden from normal Joe -- All you need to do to encode a video is to give one parameter to `ffcvt`, i.e., the path and file name of the video to be encoded, and `ffcvt` will take care of the rest, using the recommended values for both audio/video encoding to properly encode it for you.
- It can't be more simpler than that. However, beneath the simple surface, `ffcvt` is versatile and powerful enough to allow you to touch every corner of audio/video encoding. There is a huge list of environment variables which will allow you tweak the encoding methods and parameters to exactly what you prefer instead.
- Moreover, to encode a directory full of video files, including under its sub-directories, you need just to give `ffcvt` one single parameter, the directory location, and `ffcvt` will go ahead and encode all video files under that directory, including all its sub-directories as well. 

## Quick Usage

There is a quick usage help that comes with `ffcvt`, produced when it is invoked without any parameters:

```
$ ffcvt

Usage:
 ffcvt [flags] 

Flags:

  -aes  audio encoding method set (FFCVT_AES)
  -ves  video encoding method set (FFCVT_VES)
  -aea  audio encoding method append (FFCVT_AEA)
  -vea  video encoding method append (FFCVT_VEA)
  -abr  audio bitrate (64k for opus, 256k for mp3) (FFCVT_ABR)
  -crf  the CRF value: 0-51. Higher CRF gives lower quality
  (28 for x265, ~ 23 for x264) (FFCVT_CRF)

  -t    target type: x265-opus/x264-mp3 (FFCVT_T)
  -d    directory that hold input files (FFCVT_D)
  -f    input file name (either -d or -f must be specified) (FFCVT_F)
  -suf  suffix to the output file names (FFCVT_SUF)

  -ac   copy audio codec (FFCVT_AC)
  -vc   copy video codec (FFCVT_VC)
  -an   no audio, output video only (FFCVT_AN)
  -vn   no video, output audio only (FFCVT_VN)
  -vss  video: same size (FFCVT_VSS)
  -o    more options that will pass to ffmpeg program (FFCVT_O)
  -ato-opus     audio encode to opus, using -abr (FFCVT_ATO_OPUS)
  -vto-x265     video video encode to x265, using -crf (FFCVT_VTO_X265)

  -force        overwrite any existing none-empty file (FFCVT_FORCE)
  -debug        debugging level (FFCVT_DEBUG)
  -ffmpeg       ffmpeg program executable name (FFCVT_FFMPEG)

Details:

  -abr="": audio bitrate (64k for opus, 256k for mp3)
  -ac=false: copy audio codec
  -aea="": audio encoding method append
  -aes="": audio encoding method set
  -an=false: no audio, output video only
  -ato-opus=false: audio encode to opus, using -abr
  -crf="": the CRF value: 0-51. Higher CRF gives lower quality
  (28 for x265, ~ 23 for x264)
  -d="": directory that hold input files
  -debug=0: debugging level
  -f="": input file name (either -d or -f must be specified)
  -ffmpeg="ffmpeg": ffmpeg program executable name
  -force=false: overwrite any existing none-empty file
  -o="": more options that will pass to ffmpeg program
  -suf="": suffix to the output file names
  -t="x265-opus": target type: x265-opus/x264-mp3
  -vc=false: copy video codec
  -vea="": video encoding method append
  -ves="": video encoding method set
  -vn=false: no video, output audio only
  -vss=true: video: same size
  -vto-x265=false: video video encode to x265, using -crf

The `ffcvt -f testf.mp4 -debug 1 -force` will invoke

  ffmpeg -i testf.mp4 -c:a libopus -b:a 64k -c:v libx265 -x265-params crf=28 -y testf_.mkv

To use `preset`, do the following or set it in env var FFCVT_O

  cm=medium
  ffcvt -f testf.mp4 -debug 1 -force -suf $cm -- -preset $cm

Which will invoke

  ffmpeg -i testf.mp4 -c:a libopus -b:a 64k -c:v libx265 -x265-params crf=28 -y -preset medium testf_medium_.mkv

Here are the final sizes and the conversion time (in seconds):

  2916841  testf.mp4
  1807513  testf_.mkv
  1743701  testf_veryfast_.mkv   41
  2111667  testf_faster_.mkv     44
  1793216  testf_fast_.mkv       85
  1807513  testf_medium_.mkv    120
  1628502  testf_slow_.mkv      366
  1521889  testf_slower_.mkv    964
  1531154  testf_veryslow_.mkv 1413

I.e., if `preset` is not used, the default is `medium`.

Here is another set of results, sizes and the conversion time (in minutes):

 171019470  testf.avi
  55114663  testf_veryfast_.mkv  39.2
  57287586  testf_faster_.mkv    51.07
  52950504  testf_fast_.mkv     147.11
  55641838  testf_medium_.mkv   174.25

```

## Preset Method Comparison

The `ffmpeg` `x265` `preset` determines how fast the encoding process will be. if you choose `ultrafast`, the encoding process is going to run fast, but the file size will be larger when compared to `medium`. [The visual quality will be the same](https://trac.ffmpeg.org/wiki/Encode/H.265). Valid presets are `ultrafast`, `superfast`, `veryfast`, `faster`, `fast`, `medium`, `slow`, `slower`, `veryslow` and `placebo`.

Because that [the visual quality are the same](https://trac.ffmpeg.org/wiki/Encode/H.265), so there is no need to go for the slower options, because you won't be gaining anything but for the final file size. Therefore, check for yourself the above result file sizes and the conversion times, then pick a preset level you feel comfortable. The following present the same data in graphs. Click on them each to bring up bigger and most importantly, interactive graph.

[![preset small](preset-small.png)](https://fiddle.jshell.net/cL2q5p1z/3/show/ "Preset Small")

[![preset large](preset-large.png)](https://fiddle.jshell.net/nfLfd9p6/show/ "Preset Large")

I personally would go for `veryfast` because it produces the final size not much different than `fast`, `medium`, but only take less than a quarter of the time (but you may choose anything). I.e., I'll use 

    export FFCVT_O="-preset veryfast"

so as to avoid specifying it each time when invoking `ffcvt`.

## Environment Variables

For each `ffcvt` command line parameter, there is a environment variable corresponding to it. For example you can use `export FFCVT_FFMPEG=avconv` to use `avconv` instead of `ffmpeg` (Don't, I use it for my [CommandLineArgs](https://github.com/suntong001/lang/blob/master/lang/Go/src/sys/CommandLineArgs.go) to develop/test `ffcvt` without invoking `ffmpeg` each time). 

## Example: YouTube Encoding

The target type `youtube` has now been added, with settings and parameters taken from [How to Encode Videos for YouTube and other Video Sharing Sites](https://trac.ffmpeg.org/wiki/Encode/YouTube). In essence, because *"Since YouTube, Vimeo, and other similar sites will re-encode anything you give it the best practice is to provide the highest quality video that is practical for you to upload."*, every parameter has been set to aim for that high standard. I.e., a command `ffcvt -f Whatever1.mp4 -debug 1 -force -t youtube` will do:

    ffmpeg -i Whatever1.mp4 -c:a libvorbis -q:a 5 -c:v libx264 -x264-params crf=20 -pix_fmt yuv420p -y Whatever1_.mkv

All parameters are following the recommendation from the above [official ffmpeg page](https://trac.ffmpeg.org/wiki/Encode/YouTube), just that I give `crf=20` as the default because I believe it is good enough. You can still change it to `18`, which is [often considered to be roughly "visually lossless"](https://trac.ffmpeg.org/wiki/EncodingForStreamingSites), by:

    export FFCVT_CRF="18"

Of course, you can do all kinds of other tweaking as well. For a quick test, I just use the crappy video and audio:

	export FFCVT_O="-preset veryfast -q:a 2"

and it works just fine -- https://www.youtube.com/watch?v=Rms6sDp3UNM.

Note,

- The above same command, `ffcvt -f Whatever1.mp4 -debug 1 -force -t youtube` will now do:

      ffmpeg -i Whatever1.mp4 -c:a libvorbis -q:a 5 -c:v libx264 -x264-params crf=20 -pix_fmt yuv420p -y -preset veryfast -q:a 2 Whatever1_.mkv
  
- The default parameter `-q:a 5` will produce a vorbis audio of 91.9kbits/s. Overriding it with `-q:a 2`, will cause the audio to become 67.5kbits/s.
- The result video looks very crappy on YouTube. It is not because of the encoding method, as the very same file used for uploading view just fine on my machine. It looks very crappy because it has been through four rounds of encoding -- the source is from YouTube itself, with original author's and YouTube's first round encoding, plus mine, then plus the final YouTube re-encoding, the forth time, which is the final push to make the video very crappy. The 67.5kbits/s, `-q:a 2` sound however, is not bad at all. 

To recap, for high standard encoding, set 

	export FFCVT_O="-preset slow"

and optionally, as explained before:

    export FFCVT_CRF="18"

Then just use `ffcvt -f` for files or `ffcvt -d` for directories. 

## Tools Choices

As suggested before, don't use `avconv`, use `ffmpeg` instead (the `avconv` fork was more for political reasons. I personally believe `ffmpeg` is technically superior although might not be politically).

As for video/movie play back, use [mpv](http://mpv.io/). It is a fork of mplayer2 and MPlayer, and is a true *modern* *all-in-one* movie player that can play ANYTHING, and one of the few movie players being actively developed all the time. Download link is in [mpv.io](http://mpv.io/), from which Ubuntu repo I get my Ubuntu `ffmpeg` package as well. If you are unsatisfied with mpv's simple user interface, check out https://wiki.archlinux.org/index.php/Mpv#Front_ends.

