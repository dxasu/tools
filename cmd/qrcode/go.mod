module github.com/dxasu/tools/cmd/qrcode

go 1.18

replace bay.core/lancet/rain => ../../lancet/rain

require (
	bay.core/lancet/rain v0.0.0-00010101000000-000000000000
	github.com/atotto/clipboard v0.1.4
	github.com/dxasu/qrcode v1.0.0
	github.com/dxasu/tools/lancet/version v0.0.0-20230817085701-358c17504e09
	github.com/makiuchi-d/gozxing v0.1.1
	github.com/spf13/cast v1.5.1
)

require (
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
