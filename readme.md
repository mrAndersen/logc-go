Log format for Nginx:
`
log_format main_format '$http_x_real_ip - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"';

`

Available env parameters for docker container:
- CH_ADDR - Clickhouse address (default 0.0.0.0)
- CH_PORT - Clickhouse port (default 9000)
- UDP_ADDR - UDP address (default 0.0.0.0)
- UDP_PORT - UDP port (default 9222)


To start docker container execute FOR EXAMPLE
`
docker run -d --name loggo -e CH_ADDR=clickhouse -e UDP_PORT=914 -p 914:914/udp --restart always --network=my-network mrandersen7/logc-go:1.7
`

For correct work nginx should write logs to a syslog facility with format specified above. Also
make sure all ports are correctly configured and all docker networks correctly configured, you can launch docker container
in host network mode if you need to.