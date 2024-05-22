# build
```go build ipecho.go```

# configure
```
# edit backend.toml.sample
mkdir /etc/ipecho
cp backend.toml.sample /etc/ipecho/backend.toml
```

# run
```./ipecho [custom-config.toml]```

# deploy
```
# build
cp ipecho /usr/bin
cp ipecho.service /etc/systemd/system
systemctl daemon-reload
systemctl enable ipecho
systemctl start ipecho
```

# logs (when running via systemd)
```journalctl -f -u ipecho```