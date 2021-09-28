# Scanner for Neptune's Pride

<a href="https://hub.docker.com/r/albinodrought/np-scanner">
<img alt="albinodrought/np-scanner Docker Pulls" src="https://img.shields.io/docker/pulls/albinodrought/np-scanner">
</a>
<a href="https://github.com/AlbinoDrought/np-scanner/blob/master/LICENSE"><img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/np-scanner"></a>

An unofficial companion site for the game [Neptune's Pride](https://np.ironhelmet.com/) that aids you in calculating battle outcomes and collaborating with allies.

**No support will be given for this project. This is highly experimental. Breaking changes may occur at any time.**

## Screenshots

![Screenshot of Shared Map](./.readme/map-view.png)

## Features

- Share scanning data and research targets with other Neptune's Pride players
- Travel through time and view old snapshots of the universe
- Central list of all threats and their categories: red major threats that need immediate action, or green minor threats that are safe to ignore (for now!)
- Data-intensive galaxy map to find high-value target stars

## Usage

1. Add match key: `np-scanner set [game number] [key]`
2. Protect match data with code: `np-scanner protect [game number] [code]`
3. Run the Web UI on 0.0.0.0:38080: `np-scanner serve`

Other commands:

- Help: `np-scanner help`
- Test poll: `np-scanner poll [game number or "all"]`
- Disable match key: `np-scanner disable-player [game number] [player-id]`

## Building

### With Docker

`docker build -t albinodrought/np-scanner .`

### Without Docker

You need:

```
git
go
node
npm
make
```

Run `make` to generate an all-inclusive `./dist/np-scanner` binary.
