FROM golang:latest AS build
RUN git clone https://github.com/AdrienCos/pidarr_bot.git /project
WORKDIR /project
RUN CGO_ENABLED=0 go build

FROM alpine
COPY --from=build /project/pidarr_bot /app/pidarr_bot
WORKDIR /app
ENTRYPOINT [ "/app/pidarr_bot" ]