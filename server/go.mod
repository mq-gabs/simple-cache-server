module scas

go 1.26.4

require (
	go.uber.org/zap v1.28.0
	libsscas v0.0.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace libsscas => ../libs
