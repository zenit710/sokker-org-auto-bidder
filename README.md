# sokker-org-auto-bidder

## Usage

### Authorization

You need to be authorized in sokker.org to use this app. Credentials can be passed as environment variables: `SOKKER_LOGIN` for login and `SOKKER_PASS` for password.

Example:

```
SOKKER_LOGIN=player SOKKER_PASS=secret ./soab
```

### Subcommands

You need to pass subcommand as second argument to use app.

#### check-auth

Check sokker.org authorization. Example:

```
SOKKER_LOGIN=player SOKKER_PASS=secret ./soab check-auth
```

#### add

Adds new player for bid.

You need to pass flags with extra informations:

`-playerId` - player ID
`-maxPrice` - maximum price you can pay

Example:

```
SOKKER_LOGIN=player SOKKER_PASS=secret ./soab add -playerId=31235 -maxPrice=2000000
```

#### bid

Bids listed players.

Example:

```
SOKKER_LOGIN=player SOKKER_PASS=secret ./soab bid
```

### Bids automation

You should add [bid subcommand](#bid) to some scheduler like [crontab](https://man7.org/linux/man-pages/man5/crontab.5.html) to use this automation.

Example cron file entry:

```
* * * * * SOKKER_LOGIN=player SOKKER_PASS=secret /usr/bin/local/soab bid
```

## Development

You can build app after changes using `make`. `soab` executable will be built inside current directory.

## TODO

- log debug/info/warning/errors to the file: https://gist.github.com/huobazi/80805158fe1f2cac7170fe68c37741ab
- test all cases after db logic change
- split tools into directories
- add workflows to: unit test, lint test, PR approve, release build
