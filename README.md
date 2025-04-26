# fr24-kml-splitter

This is a very small tool that can edit flightradar24 KML track. It does two things:

- Removes the _trail_ lines from the KML file. FR24 adds both the recorded positions and trail lines to the files, sometimes you just want the positions without the trail.
- Optionally splits the track at the timestamp you specify.

## Usage

### General

Specify the track file as first argument:

```shell
go run main.go ~/Downloads/track.kml
```

The modified file will be written to this directory as `route-only.kml`

### Splitting the track

To split the track, you need to specify the first timestamp of the new flight as second argument:

```shell
go run main.go ~/Downloads/track.kml 2025-04-08T15:40:35Z
```

This will split the track before this timestamp and create two distinct track files: `route-only-first.kml` and `route-only-second.kml`

## Contributing

You need to have pre-commit installed. Run `pre-commit install` before commiting.

This project uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).
