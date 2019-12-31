templateFile=$GOPATH/src/github.com/go-easygen/easygen/test/commandlineFlag
[ -s $templateFile.tmpl ] || templateFile=/usr/share/gocode/src/github.com/go-easygen/easygen/test/commandlineFlag
[ -s $templateFile.tmpl ] || {
  echo No template file found
  exit 1
}

easygen $templateFile ffcvt_cli | sed '/\tVES\t\tstring/{ N; N; N; N; N; N; N; N; N; s|^.*$|\tEncoding\t// anonymous field to hold encoding values|; }; /\tExt\t\tstring/d; ' | gofmt > config.go
