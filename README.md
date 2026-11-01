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

### Finding a stop ID

Open the [HSL GraphiQL explorer](https://api.digitransit.fi/graphiql/hsl), paste this query
and run it (Cmd/Ctrl+Enter):

```graphql
{
  stops(name: "Kamppi") {
    gtfsId
    name
    code
  }
}
```

A station usually has several stops, one per platform or direction. To see which direction
each one goes, query their next departures (the `headsign` is the direction):

```graphql
{
  stops(ids: ["HSL:1040601", "HSL:1040602"]) {
    gtfsId
    name
    stoptimesWithoutPatterns(numberOfDepartures: 3) {
      serviceDay
      realtimeArrival
      realtime
      headsign
    }
  }
}
```

Put the `gtfsId` values of the stops you want in `ARRIVAL_IDS`.

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
