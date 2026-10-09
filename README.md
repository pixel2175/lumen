# lumen

Universal brightness controller for Linux. No external dependencies.

## Install

Requires Go.

```sh
git clone https://github.com/pixel2175/lumen
cd lumen
make install
```

```sh
make uninstall
make clean
```

## Usage

```sh
lumen set 50                  # set to 50%
lumen set +10%                # increase by 10
lumen set -5                  # decrease by 5
lumen get                     # current brightness (%)
lumen get monitors            # list monitors
lumen -m DP-1 set 70          # choose a monitor
lumen -m=DP-1 get
```

## Permissions

### Laptop panel (sysfs)

Writing to `/sys/class/backlight/*/brightness` needs root by default. Add a udev rule to allow the `video` group:

```sh
echo 'ACTION=="add", SUBSYSTEM=="backlight", RUN+="/bin/chgrp video $sys$devpath/brightness", RUN+="/bin/chmod g+w $sys$devpath/brightness"' | sudo tee /etc/udev/rules.d/90-backlight.rules
sudo usermod -aG video $USER
```

### External monitors (DDC/CI)

Load the `i2c-dev` kernel module:

```sh
sudo modprobe i2c-dev
echo i2c-dev | sudo tee /etc/modules-load.d/i2c.conf   # load on boot
```

Allow your user to access `/dev/i2c-*`:

```sh
sudo groupadd -f i2c
sudo usermod -aG i2c $USER
echo 'KERNEL=="i2c-[0-9]*", GROUP="i2c", MODE="0660"' | sudo tee /etc/udev/rules.d/99-i2c.rules
sudo udevadm control --reload-rules && sudo udevadm trigger
```
Also make sure DDC/CI is enabled in your monitor's on-screen menu.
