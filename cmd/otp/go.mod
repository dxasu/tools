module github.com/dxasu/tools/cmd/otp

go 1.18

replace bay.core/lancet/rain => ../../lancet/rain

require (
	bay.core/lancet/rain v0.0.0-00010101000000-000000000000
	github.com/atotto/clipboard v0.1.4
	github.com/dxasu/tools/lancet/version v0.0.0-20230817102921-3ddac65cae69
	github.com/pquerna/otp v1.4.0
)

require github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
