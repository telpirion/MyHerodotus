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

## Change Python version on Cloud Shell

Use `pyenv` to switch Python versions.

1. Install `pyenv` to install python on persistent home directory.

    ```sh
    $ curl https://pyenv.run | bash
    ```

1. Add `pyenv` to path.

    ```sh
    $ echo 'export PATH="$HOME/.pyenv/bin:$PATH"' >> ~/.bashrc
    $ echo 'eval "$(pyenv init -)"' >> ~/.bashrc
    $ echo 'eval "$(pyenv virtualenv-init -)"' >> ~/.bashrc
    ```

1. Update `.bashrc`.

    ```sh
    $ source ~/.bashrc
    ```

1. Install desired version of Python and make default.

    ```sh
    $ pyenv install 3.10.7
    $ pyenv global 3.10.7
    ```

1. Switch to the desired version of Python for `virtualenv`.

    ```sh
    virtualenv env --python=python3.10.7
    ```