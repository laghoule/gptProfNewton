FROM golang:1.22-alpine AS dep
WORKDIR /src/
COPY . .
RUN cd cmd \
	go get -d -v

FROM dep AS build
ARG VERSION "devel"
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.version=$VERSION'" -o gptProfNewton cmd/main.go cmd/printer.go

FROM alpine:3.19
LABEL org.opencontainers.image.source https://github.com/laghoule/gptProfNewton
COPY --from=build /src/gptProfNewton /usr/bin/
USER nobody
ENV OPENAI_API_KEY ""
ENTRYPOINT ["/usr/bin/gptProfNewton"]
