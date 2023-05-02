module github.com/valyala/quicktemplate/examples/basicserver

go 1.20

replace github.com/valyala/quicktemplate v0.0.0 => ../../

require (
	github.com/valyala/fasthttp v1.47.0
	github.com/valyala/quicktemplate v0.0.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)
