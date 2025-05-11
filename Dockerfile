### build stage
FROM golang:1.23.9-alpine3.21 AS build


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN apk add protoc build-base

# также нужно в стадии билда создать этого юзера и на всякий случай группу
RUN addgroup -S nonroot
RUN adduser -S nonroot -G nonroot

# RUN export PATH="$PATH:$(go env GOPATH)/bin" <-- Неправильно. Правильно ниже
ENV PATH="$PATH:$(go env GOPATH)/bin"
ENV CGO_ENABLED=1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY ./ .


# to generate go files from .proto
RUN protoc -I proto proto/auth.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

# compile application during build rather than at runtime
# add flags to statically link library
RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o go_jwt_mcs ./cmd




### Runtime stage
FROM scratch

# без этого не будет работать, выведется ошибка:
# 
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/go_jwt_mcs go_jwt_mcs

USER nonroot

EXPOSE 50051

CMD ["/go_jwt_mcs"]