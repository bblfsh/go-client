module github.com/bblfsh/go-client/v4

go 1.12

require (
	github.com/bblfsh/sdk/v3 v3.2.5
	github.com/jessevdk/go-flags v1.4.0
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.20.1
	gopkg.in/bblfsh/sdk.v1 v1.17.0
	gopkg.in/src-d/go-errors.v1 v1.0.0
)

replace github.com/bblfsh/sdk/v3 v3.2.5 => github.com/lwsanty/sdk/v3 v3.2.1-0.20191112180841-f3c9d3e2f493
