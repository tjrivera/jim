jim
===

Jim is a simple service monitor that texts a number when a service is down

It relies on the Twilio API to send messages to a defined phone number.

Make sure to set these variables in your environement

`TWILIO_SID`
`TWILIO_TOKEN`
`SMS_FROM`
`SMS_TO`

```bash
$ docker build -t jim .
$ docker run \
    -e TWILIO_SID=ABC \
    -e TWILIO_TOKEN=123 \
    -e SMS_FROM=555555555 \
    -e SMS_TO=4444444444 \
    -p 8080:8080 jim /go/bin/jim -target="http://example.com"
```

you can also pass `-poll=[n]s` to specify the number of seconds between polling.

Jim will send you a text after 5 failed attempts on the target host.
