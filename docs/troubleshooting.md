# Troubleshooting

## Clear space on Cloud Shell

1. Get the usage of local files.

    ```sh
    $ df -h
    $ du -hs $(ls -A)
    ```

1. Clean pip cache.

    ```sh
    $ pip cache purge
    ```

1. Clean the Go cache.

    ```sh
    $ go clean --modcache
    $ go clean --cache
    ```

1. Prune the system (Docker).

    ```sh
    $ docker system prune
    ```

1. Run Git garbage collection

    ```sh
    $ git gc
    ```

1. Delete npm.

    ```sh
    $ rm -rfd .npm
    ```
