With run this script I use a tool that called [HTTPie](https://httpie.org/)

At first you should install it or use some thing other instead it(you need write a script yourself).
`pip install httpie`

Note:
you need check the script as your server Settings.(the `step` value in the script)

Usage:

1. set your config in config.txt. 
2. run `./run_homeip.sh`

config.txt like:
`https://[domain.com/myip] [yourkey] [step] [auth:auth]`

Crontab:
`*/5 * * * * cd [/your/path/] && ./run_homeip.sh >> response.log 2>&1`
