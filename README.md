# komorebit
Yet another tray indicator for [komorebi](https://github.com/LGUG2Z/komorebi/)

![image](https://github.com/user-attachments/assets/7693ece1-82bd-4c44-bf2d-f4149995e19d)

![image](https://github.com/user-attachments/assets/192e76ba-f3ad-4280-a109-0bfd367f95cb)

## Features
- Active workspace indicator (active monitor) and switching menu
- Acive layout inidicator (active monitor) and switching menu
- Komorebi pause indicator and button
- Komorebi restart button (komorebic stop; komorebic start)

## Similar tools / inspired by
### Compared to [komotray](https://github.com/joshprk/komotray)
- Doesn't wrap around komorebi
- Doesn't come with keybindings
### Compared to [komotray](https://github.com/joshprk/komotray)
- Has a menu with a few buttons

## Installation
For daily use, make sure to build with a console-hiding flag:
```
go build -ldflags "-H=windowsgui"
```

## Known issues
- No active monitor indication
- Workspace switching menu only affects default monitor and always activates it
