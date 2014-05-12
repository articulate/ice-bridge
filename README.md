# ice-bridge

Dropbox Archiver written in Go

## Getting Started

### Configuration

Install Go and setup your [GOPATH](http://golang.org/doc/code.html#GOPATH).

Once your go environment is up and running the application itself requires some
configuration. It will require you to create a dropbox application at their app
center. Once that is created you will need the CLIENT_ID and CLIENT_SECRET.

By default the application looks for it's configuration file in your home
directory. On OSX it will be `$HOME/.icebridge`.

The configuration file is a JSON formatted file with the following structure.

```json
{
  "client_id": "APP_CLIENT_ID",
  "client_secret": "APP_CLIENT_SECRET",
  "token": "YOUR OAUTH TOKEN",
  "local_path": "FULL PATH TO ARCHIVE FILES LOCALLY",
  "dropbox_path": "PATH ON DROPBOX TO ARCHIVE"
}
```

If you do not have a token, the application will walk you through the steps
required to authorize the Dropbox Application and obtain an access token.

### Running

With a working Go environment and the configuration file written you can
compile icebridge on your machine using.

```shell
go build
```

Then you can run it from the current directory using
```shell
./ice-bridge
```

You can also compile the binary and place it in your `$GOPATH/bin` using

```shell
go install
```

Then you can simply type
```shell
ice-bridge
```
to run the application.

