FROM golang:1.18-alpine as build
RUN apk upgrade --no-cache --force
RUN apk add --update build-base make git
WORKDIR /go/src/github.com/vibin18/go-shares

# Compile
COPY ./ /go/src/github.com/vibin18/go-shares
RUN make dependencies
RUN make build


# Final Image
FROM gcr.io/distroless/static AS export-stage



COPY --from=build /go/src/github.com/vibin18/go-shares/go-shares /
COPY --from=build /go/src/github.com/vibin18/go-shares/templates/* /templates/
USER 1000:1000
ENTRYPOINT ["/go-shares"]