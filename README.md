# Scanner for Neptune's Pride

<p align="center">
  <img src="./.readme/banner.png">
  <p align="center">
    <a href="https://github.com/AlbinoDrought/np-scanner/blob/master/LICENSE"><img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/np-scanner"></a>
  </p>
</p>

An unofficial companion site for the game [Neptune's Pride](https://np.ironhelmet.com/) that aids you in calculating battle outcomes and collaborating with allies.

**No support will be given for this project. This is highly experimental. Breaking changes may occur at any time.**

## Screenshots

![Screenshot of Shared Map](./.readme/map-view.png)

## Features

- Share scanning data and research targets with other Neptune's Pride players
- Travel through time and view old snapshots of the universe
- Central list of all threats and their categories: red major threats that need immediate action, or green minor threats that are safe to ignore (for now!)
- Early warning system alerts you when a star is targetted
- Data-intensive galaxy map to find high-value target stars

## Usage

1. Add match key: `np-scanner set [game number] [key]`
2. Protect match data with code: `np-scanner protect [game number] [code]`
3. Run the Web UI on 0.0.0.0:38080: `np-scanner serve`

Other commands:

- Help: `np-scanner help`
- Test poll: `np-scanner poll [game number or "all"]`
- Disable match key: `np-scanner disable-player [game number] [player-id]`
- Create a code that can only see data from player 1 and 2: `np-scanner protect --allowed-uid 1 --allowed-uid 2 [game number] [code]`
- Replace all other codes: `np-scanner protect --wipe [game number] [code]`
- Associate a game player with their Discord user ID for notifications: `np-scanner set-discord [game number] [player uid] [discord user id]`

Config:

- Discord Webhook URL for alerts: env var `NP_SCANNER_DISCORD_WEBHOOK_URL=https://...` or cli arg `--discord-webhook-url=https://...`

## Building

### With Docker

`docker build -t albinodrought/np-scanner .`

#### Sample Docker Compose

```yml
version: '2'
services:
  np-scanner:
    image: albinodrought/np-scanner
    volumes:
    - /some/local/path:/data
    command:
    - /np-scanner
    - --db-path=/data/np.db
    - serve
    ports:
      - 38080:38080
```

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


