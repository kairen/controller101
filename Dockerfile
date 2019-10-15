FROM kairen/golang-dep:1.12-alpine AS build
LABEL maintainer="Kyle Bai <kyle.b@k2r2bai.com>"

ENV GOPATH "/go"
ENV PROJECT_PATH "$GOPATH/src/github.com/cloud-native-taiwan/controller101"
ENV GO111MODULE "on"

COPY . $PROJECT_PATH
RUN cd $PROJECT_PATH && \
  make && mv out/controller /tmp/controller

# Running stage
FROM alpine:3.7
COPY --from=build /tmp/controller /bin/controller
ENTRYPOINT ["controller"]