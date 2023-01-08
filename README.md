# Waybar Unread Email

A simple Go program to print unread email counts in JSON suitable to use as
[Waybar custom module](https://man.archlinux.org/man/waybar-custom.5.en).

![image](https://user-images.githubusercontent.com/10193999/161427391-244b302a-5bea-4ed5-88c6-42eace75f568.png)

## Install

Build and copy `waybar-unread-email` to somewhere in the `PATH`.

```txt
go build .
sudo cp ./waybar-unread-email /usr/local/bin
```

## Usage

```sh
Usage of waybar-unread-email:
  -config string
        path to YAML file containing configuration (default "/home/user/.config/waybar-unread-email/config.yaml")
  -output string
        output format, one of: "json", "yaml", "num" (default "json")
```

## Example Output

### JSON (default)

```json
// 3 unread messages
{
  "text": "3",
  "tooltip": "Proton: 0 unread\nGMail: 2 unread\nOutlook: 1 unread",
  "class": "unread",
  "percentage": 100
}
```

```json
// 0 unread messages, showZero: false
{
  "text": "",
  "tooltip": "Proton: 0 unread\nGMail: 0 unread\nOutlook: 0 unread",
  "percentage": 0
}
```

```json
// 0 unread messages, showZero: true
{
  "text": "0",
  "tooltip": "Proton: 0 unread\nGMail: 0 unread\nOutlook: 0 unread",
  "percentage": 0
}
```

### YAML

```yaml
class: unread
percentage: 100
text: "3"
tooltip: |-
  Proton: 0 unread
  GMail: 2 unread
  Outlook: 1 unread
```

### Num

```txt
3
```

## Configuration

Configuration is stored in a YAML file that contains settings for a list IMAP
servers to be checked for unread messages. If `encryption` is set to `TLS` or
`STARTTLS` encryption will be attempted. To ignore invalid certificates set
`skipVerify` to `true`. By default the output will be empty if there are no
unread messages. To have `0` printed set `showZero` to `true`.

```yaml
showZero: false
servers:
  - name: Proton
    address: localhost:1143
    username: road.runner@pm.me
    password: isecretlylovewile
    security: STARTTLS
    skipVerify: true
  - name: GMail
    address: imap.gmail.com:993
    username: wile.e.coyote@gmail.com
    password: RocketShoes!!
    security: TLS
    skipVerify: false
  - name: Outlook
    address: outlook.office365.com:993
    username: bugs.bunnye@outlook.com
    password: password123
    security: TLS
    skipVerify: false
```

## Waybar config

```json
{
  "layer": "top",
  "position": "top",
  "modules-left": ["custom/waybar-unread-email"],

  "custom/waybar-unread-email": {
    "exec": "waybar-unread-email",
    "return-type": "json",
    "interval": 300,
    "on-click": "geary",
    "on-click-right": "waybar-unread-email",
    "format": "{icon}{}",
    "format-icons": ["", "﫮"]
  }
}
```

## Libraries

- [go-imap](https://github.com/emersion/go-imap)
- [sigs.k8s.io/yaml](https://github.com/kubernetes-sigs/yaml)
