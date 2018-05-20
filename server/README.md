you can set as a server service with Systemd.
also run script to start or stop it by manually.

Note:
this scripts should be in the same directory as the homeip executable file.

Systemd:
place your config as `/lib/systemd/system/homeipkit.service`
```
sudo systemctl daemon-reload

sudo systemctl enable homeipkit

sudo systemctl start homeipkit

sudo systemctl status homeipkit.service

sudo systemctl stop homeipkit
```