# pidarr_bot - A telegram bot to control Pidarr

This telegram bot is designed to be used with (Pidarr)[https://github.com/AdrienCos/pidarr].


To build it, run `make` from the root of the repo.

To run it, start by settings the following environment variables:

| Variable         | Contents                                                                                                       | Default value    |
| ---------------- | -------------------------------------------------------------------------------------------------------------- | ---------------- |
| `PIDARR_TOKEN`   | Telegram bot token                                                                                             | None             |
| `PIDARR_CHATID`  | ID of the chat between the bot and yourself                                                                    | None             |
| `RADARR_TOKEN`   | Radarr API token                                                                                               | None             |
| `RADARR_HOST`    | Host and path to your Radarr instance                                                                          | `localhost:7878` |
| `RADARR_PATH`    | Path where Radar should add newly added movies                                                                 | `/movies`        |
| `RADARR_QUALITY` | ID of the quality profile to use for newly added movies (query the `/profile` Radarr endpoint to get them all) | `4` (HD 1080p)   |

You can then run the bot with `./pidarr_bot`.

## Using Docker

This project comes with a Dockerfile to build the bot in a dedicated image. To build it, simply run `docker build -t $IMAGE_NAME:$IMAGE_TAG .` from the root of the repo. 

To run the container, first create a environment file containing the variables listed above, with the following structure:

```
PIDARR_TOKEN=$telegram_bot_token
PIDARR_CHATID=$telegram_chat_id
RADARR_TOKEN=$radarr_api_token
RADARR_HOST=$radarr_host
RADARR_PATH=$radarr_path
RADARR_QUALITY=$radarr_quality
```

You can then create a new container from your image and run it with:

```
docker run --env-file $ENV_FILENAME $IMAGE_NAME:$IMAGE_TAG
```

# How it works

Once the bot is running, you can ask it to search for movies by sending it a message in the form `/movies Movie Name`. The bot will then query the Radarr API to find a list of all movies that match that search term, and return you a list of options in the form an inline keyboard. 

Find the movie(s) you want in the list and simply click them. The bot will then call Radarr again to have it add and search for the movies.