#! /bin/sh

../ffcvt -version

# Test -sym control
../ffcvt -t x265-opus -n -d . > /tmp/ffcvt_test.txt 2>&1
../ffcvt -t x265-opus -n -d . -sym  >> /tmp/ffcvt_test.txt 2>&1

../ffcvt -n -d .  >> /tmp/ffcvt_test.txt 2>&1
../ffcvt -n -d . -sym  >> /tmp/ffcvt_test.txt 2>&1

../ffcvt -n -sym -debug 2 -d . -w /tmp >> /tmp/ffcvt_test.txt 2>&1

sed -i '/ [0-9.]*[mµ]*s$/d' /tmp/ffcvt_test.txt
diff -wU 1 ffcvt_test.txt /tmp/ffcvt_test.txt

ret=$?
echo $ret
exit $ret
