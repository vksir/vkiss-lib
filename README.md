# Vkiss Lib


## Bootstrap

1. Get the config path
2. Read the config file
   1. Read one
   2. Write default config, then read it
3. Init workspace
4. Init Log
5. Init database
6. Run


## Build

```bash
export GOOS=linux 
export GOARCH=amd64
go build 
```

## Run

### Ddns

```bash
./vkiss ddns install server
```

```bash
./vkiss ddns install monitor
```
