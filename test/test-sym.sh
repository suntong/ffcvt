#! /bin/sh

# Test -sym control
../ffcvt -n -d . -debug 1 > /tmp/ffcvt_test.txt 2>&1 
../ffcvt -n -d . -debug 1 -sym  >> /tmp/ffcvt_test.txt 2>&1 

sed -i '/ [0-9.]*[mµ]*s$/d' /tmp/ffcvt_test.txt
diff -wU 1 /tmp/ffcvt_test.txt .

ret=$?
echo $ret
exit $ret
