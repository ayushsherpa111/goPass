# goPass
CLI password manager.

The gopass_d folder contains a daemon that is meant to run in the background.
Place the `gopass.service` unit file in `/usr/lib/systemd/system` 

Before starting the daemon, 
- Compile the source code 
```bash
cd gopass_d && go build
```
- Place the binary in `/bin` 
```bash
sudo mv goPassd /bin
```

- Start the daemon with
```bash
sudo systemctl start gopass
```
- Compile the client
```bash
cd goPass_client && go build 
```

# Commands
`goPass --help`

```bash
goPass is a password manager that runs entirely on your CLI

Usage:
  goPass [flags]
  goPass [command]

Available Commands:
  add         Add a new password to your vault
  delete      Delete a 
  get         Get your passwords [Decrypted].
  help        Help about any command
  init        Initialize your Repo

Flags:
  -h, --help   help for goPass

Use "goPass [command] --help" for more information about a command.
```



