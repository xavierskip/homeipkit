With run this script I use a tool that called [HTTPie](https://httpie.org/)

At first you should install it or use some thing other instead it(you need write a script yourself).
`pip install httpie`

Note:
you need check the script as your server Settings.(the `step` value in the script)

Usage:

`./myip.sh "https://domain.com/myip" "yourkey"`


Crontab:

`*/5 * * * * /your/path/to/myip.sh "https://domain.com/myip" "yourkey">> /your/path/to/response.log 2>&1`