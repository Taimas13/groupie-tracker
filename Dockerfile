FROM golang:alpine AS appbuild

WORKDIR /g-t-docker
COPY . /g-t-docker/

RUN go build -o main cmd/main.go

FROM alpine

LABEL maintainer="TaimasAndAliser-99"

WORKDIR /g-t-app-dir
COPY --from=appbuild /g-t-docker/main /g-t-app-dir/
COPY --from=appbuild /g-t-docker/ui/ /g-t-app-dir/ui
COPY --from=appbuild /g-t-docker/config/config.json /g-t-app-dir/config/

CMD [ "./main" ]