# Troubleshooting

## Clear space on Cloud Shell

1. Get the usage of local files.

```sh
$ df -h
$ du -hs $(ls -A)
```

2. Clean the Go cache.

```sh
$ go clean --cache
```

3. Prune the system (Docker).

```sh
$ docker system prune
```

4. Run Git garbage collection

```sh
$ git gc
```

5. Delete npm

```sh
$ rm -rfd .npm
```