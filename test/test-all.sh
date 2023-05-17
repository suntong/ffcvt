#! /bin/sh

#env

FFCVT=../ffcvt
# During Debian building, it is at obj-x86_64-linux-gnu/bin/ffcvt
# while the source is at obj-x86_64-linux-gnu/src/github.com/suntong/ffcvt/test
[ -s $FFCVT ] || FFCVT=../../../../../bin/ffcvt

$FFCVT -version

echo
echo '# Test (config.go) cli help output'
$FFCVT > /tmp/ffcvt_test.txt 2>&1

echo '# Test transcoding single file' | tee -a /tmp/ffcvt_test.txt
$FFCVT -n -debug 0 -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1

echo '# Test transcoding different target types' | tee -a /tmp/ffcvt_test.txt
$FFCVT -t webm -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t x265-opus -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t x264-mp3 -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t wx -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t youtube -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t copy -n -f StreamSample.mkv -w /tmp >> /tmp/ffcvt_test.txt 2>&1


echo '# Test -sym control' | tee -a /tmp/ffcvt_test.txt
$FFCVT -t x265-opus -n -d . >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -t x265-opus -n -d . -sym  >> /tmp/ffcvt_test.txt 2>&1

$FFCVT -n -d .  >> /tmp/ffcvt_test.txt 2>&1
$FFCVT -n -d . -sym  >> /tmp/ffcvt_test.txt 2>&1

$FFCVT -n -sym -debug 2 -d . -w /tmp >> /tmp/ffcvt_test.txt 2>&1

echo '# Compare test results, 0 means AOK:'
sed -i '/ [0-9.]*[nmµ]*s$/s// xxx ms/' /tmp/ffcvt_test.txt
diff -wU 1 ffcvt_test.txt /tmp/ffcvt_test.txt

ret=$?
echo $ret
exit $ret
