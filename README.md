messenger
=========

little messenger tool to send messages either to libnotify, or via tcp to someone else. see also messaged.

Use it
------
display output of find via libnotify. useful for commands that take their time to display the result as notification when its finished:
```bash
find /bin -iname "b*" | messenger
```

use messaged to send messages SSL encrypted via tcp - can of course be used to send messages to another pc:
```bash
messaged &
messenger --type tcp --title Hello --message World --to localhost
```

Get it
-------
```bash
go get github.com/ProhtMeyhet/messenger
```

Dependencies
-------------
libnotify

https://github.com/ProhtMeyhet/libgomessage

Licence
-------
see LICENCE-AGPL
