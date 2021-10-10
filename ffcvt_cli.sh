templateFile=$GOPATH/src/github.com/go-easygen/easygen/test/commandlineFlag
[ -s $templateFile.tmpl ] || templateFile=/usr/share/gocode/src/github.com/go-easygen/easygen/test/commandlineFlag
[ -s $templateFile.tmpl ] || templateFile=/usr/share/doc/easygen/examples/commandlineFlag
[ -s $templateFile.tmpl ] || {
  echo No template file found
  exit 1
}

easygen $templateFile ffcvt_cli | tee /tmp/ffcvt_cli.go | sed '/\tVES\t\tstring/{ N; N; N; N; N; N; N; N; N; s|^.*$|\tEncoding\t// anonymous field to hold encoding values|; }; /\tExt\t\tstring/d; /flag.MFlagsVar/{s/flag.MFlagsVar/flag.Var/; s/, "",/,/; } ' | gofmt > config.go
