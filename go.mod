module github.com/bukodi/demo-app

go 1.24

require (
	github.com/aws/aws-lambda-go v1.49.0
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.15
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.19.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.80.1
	github.com/awslabs/aws-lambda-go-api-proxy v0.16.2
	github.com/containerd/errdefs v1.0.0
	github.com/docker/docker v28.2.2+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/glebarez/sqlite v1.11.0
	github.com/go-oauth2/oauth2/v4 v4.5.3
	github.com/go-session/session v3.1.2+incompatible
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/kr/pretty v0.3.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.9.1
	github.com/tebeka/selenium v0.9.9
	golang.org/x/crypto v0.39.0
	golang.org/x/net v0.41.0
	golang.org/x/oauth2 v0.30.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.0
	gotest.tools/v3 v3.5.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.68 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.25.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.20 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/sys/atomicwriter v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tidwall/btree v0.0.0-20191029221954-400434d76274 // indirect
	github.com/tidwall/buntdb v1.1.2 // indirect
	github.com/tidwall/gjson v1.14.2 // indirect
	github.com/tidwall/grect v0.0.0-20161006141115-ba9a043346eb // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/rtree v0.0.0-20180113144539-6cd427091e0e // indirect
	github.com/tidwall/tinyqueue v0.0.0-20180302190814-1e39f5511563 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	modernc.org/libc v1.61.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.8.2 // indirect
	modernc.org/sqlite v1.36.1 // indirect
)

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1
