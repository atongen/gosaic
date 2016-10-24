FROM golang:1.7

WORKDIR /home/app
RUN mkdir src pkg bin vendor
COPY ./src/ /home/app/src/
COPY ./vendor/ /home/app/vendor/
COPY ./Makefile /home/app/
COPY ./version /home/app/
RUN make

ENTRYPOINT ["/home/app/bin/gosaic"]
