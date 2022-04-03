# unread-email
A simple script to print unread emails count and suitable to use as Waybar module

![image](https://user-images.githubusercontent.com/10193999/161427391-244b302a-5bea-4ed5-88c6-42eace75f568.png)

## Usage
```
Usage of ./unread-count:
  -P string
    	Password
  -S string
    	Server
  -U string
    	Username
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
        "exec": "unread-count -U user -P pass -S 127.0.0.1:1143",
        "return-type": "json",
        "interval": 30,
        "format": "{}",
        "on-click": "geary"
    },
}
```

