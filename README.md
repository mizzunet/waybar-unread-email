# unread-email
A simple script to print unread emails count and suitable to use as Waybar module

![image](https://user-images.githubusercontent.com/10193999/161427391-244b302a-5bea-4ed5-88c6-42eace75f568.png)

## Usage
```
usage: unread-count [-h] -u USERNAME -p PASSWORD -S SERVER -P PORT

options:
  -h, --help            show this help message and exit
  -u USERNAME, --username USERNAME
  -p PASSWORD, --password PASSWORD
  -S SERVER, --server SERVER
  -P PORT, --port PORT

```

## Install
`python-imaplib2` required as a dependency.

Copy `unread-count` to `PATH`

## Waybar config
```
{
"layer": "top",
"position": "top", 
"modules-left": ["custom/unreadcount"],

    "custom/unreadcount": {
        "exec": "unread-count -u user -p pass -S 127.0.0.1 -P 1143",
        "return-type": "json",
        "interval": 30,
        "format": "{}",
        "on-click": "geary"
    },
}
```

