> a learning note, refered from https://eddycjy.com/posts/prometheus/2020-05-16-metrics/

## Intro
To present how to build one Prometheus server and another client server, then build a metrics transmission between them.

## How to use
0. download Prometheus
    - download from https://prometheus.io/download/
    - use `tar xvfz prometheus-*.tar.gz` and mv outputs to prometheus-server, just exclude `prometheus.yml`
    - config prometheus by changing the file `prometheus.yml`
1. start prometheus server
    ```
    $ cd prometheus-server
    $ ./prometheus # it will run server at 9090 port as default, you can change it using --web.listen-address="0.0.0.0:9091"
    ```
2. start client service
    ```
    $ cd client-service
    $ go run .
    ```
3. check metrics
    - type `localhost:9090` in browser to visit the web GUI of prometheus, then you can type any go MemStats params to query the result
    - type `localhost:10001/metrics` in browser to fetch all metrics from the client service
4. manipulate metrics
    - `localhost:10001/counter`, then query `api_requests_total` to find the change
    - `localhost:10001/gauge?num=10`, then query `queue_num_total` to find the change
    - `localhost:10001/histogram`, then query `http_durations_histogram_seconds` to find the change
    - `localhost:10001/summary`, then query `http_durations_seconds` to find the change

