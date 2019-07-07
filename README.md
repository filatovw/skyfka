# Skyfka

This tool allows you to see images in Skype.


# Usage

Download one of packages for your system

Application requires read+write permissions to your `hosts` file.

One-time execution:

    sudo ./bin/skyfka
    2019/06/26 20:33:46 skyfka started
    2019/06/26 20:33:46 patch host: api.asm.skype.com
    2019/06/26 20:33:46 skyfka stopped

Regular execution:

    sudo ./bin/skyfka -regular
    2019/06/26 20:34:17 skyfka started
    2019/06/26 20:34:17 monitor api.asm.skype.com
    ^C2019/06/26 20:34:20 stopped with signal: interrupt
    2019/06/26 20:34:20 skyfka stopped

# TODO:

- add PID file
- add configurations for Systemd, supervisord
- build in Docker

# Development

Lint code

    make lint

Build current version

    make build

Delete old artifacts

    make clean

Build release for all available platforms

    make release
