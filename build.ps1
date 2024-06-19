$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -tags lambda.norpc -o build/bootstrap cmd/main.go
~\Go\Bin\build-lambda-zip.exe -o build/myFunction.zip build/bootstrap