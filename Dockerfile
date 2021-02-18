FROM golang:1.15
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN make bootstrap && make build
CMD ["make run"]