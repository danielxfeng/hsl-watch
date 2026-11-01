# hsl-watch

A tiny Go server that polls the HSL (Helsinki public transport) API every 10 seconds
and serves the next arrivals for your stops as JSON. Small enough to run on a Raspberry Pi.

## Setup

1. Get an API key from the [Digitransit developer portal](https://digitransit.fi/en/developers/) (free).
2. Copy the sample config and fill it in:

   ```sh
   cp .env.sample .env
   ```

   | Variable       | Meaning                                              |
   | -------------- | ---------------------------------------------------- |
   | `HSL_API_KEY`  | Your Digitransit key (required)                      |
   | `ARRIVAL_IDS`  | Comma-separated stop IDs, e.g. `HSL:1040601,HSL:1040602` |
   | `DEBUG`        | `true` to log the full payload on every fetch        |

## Run

```sh
make run                    # local, needs Go
docker compose up -d --build   # Docker (works on a Pi too)
```

## API

| Endpoint               | Description                                              |
| ---------------------- | -------------------------------------------------------- |
| `GET /healthz`         | Returns `{"status":"ok"}` while the server is up         |
| `GET /api/v1/arrivals` | Cached arrivals; `503` until the first successful fetch  |

Arrival time as a unix timestamp = `serviceDay + realtimeArrival`.
`fetchedAt` is when the data was last fetched, so clients can tell if it is stale.

To try the API from VS Code, install the REST Client extension and open `api.http`.
