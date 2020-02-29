FROM golang:latest AS build
RUN git clone https://github.com/AdrienCos/pidarr_bot.git
WORKDIR /go/pidarr_bot
RUN go build

FROM golang:latest
COPY --from=build /go/pidarr_bot/pidarr_bot /pidarr_bot/pidarr_bot
CMD [ "/pidarr_bot/pidarr_bot" ]