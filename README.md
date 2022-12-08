# coins-exporter

![Latest release](https://img.shields.io/github/v/release/s16rv/coins-exporter)
[![Actions Status](https://github.com/s16rv/coins-exporter/workflows/test/badge.svg)](https://github.com/s16rv/coins-exporter/actions)

coins-exporter is a Prometheus scraper that fetches the data from Coingecko API.

## What can I use it for?

You can run coins-exporter to scrape the data from it and you can access coins exporter from listen address such as below:
```
http://127.0.0.1:9500/metrics/coins?ids=juno-network,cosmos,osmosis
```
Then, the response would be like this:
```
# HELP coins_price Price of Coins in currency
# TYPE coins_price gauge
coins_price{currency="USD",id="cosmos",name="Cosmos Hub",symbol="atom"} 9.59
coins_price{currency="USD",id="juno-network",name="JUNO",symbol="juno"} 1.7
coins_price{currency="USD",id="osmosis",name="Osmosis",symbol="osmo"} 0.909976
```

## How can I set it up?

First of all, you need to download the latest release from [the releases page](https://github.com/s16/coins-exporter/releases/). After that, you should unzip it and you are ready to go:

```sh
wget <the link from the releases page>
tar xvfz coins-exporter-*
./coins-exporter
```

That's not really interesting, what you probably want to do is to have it running in the background. For that, first of all, we have to copy the file to the system apps folder:

```sh
sudo cp ./coins-exporter /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/coins-exporter.service
```

You can use this template (change the user to whatever user you want this to be executed from. It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=Cosmos Exporter
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=coins-exporter
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Then we'll add this service to the autostart and run it:

```sh
sudo systemctl enable coins-exporter
sudo systemctl start coins-exporter
sudo systemctl status coins-exporter # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u coins-exporter -f --output cat
```

## How can I scrape data from it?

Here's the example of the Prometheus config you can use for scraping data:

```yaml
scrape-configs:
  # specific coin(s)
  - job_name: coins-exporter
    scrape_interval: 15s
    metrics_path: /metrics/coins
    static_configs:
      - targets:
        - juno-network,cosmos,osmosis
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_ids
      - source_labels: [__param_ids]
        target_label: instance
      - target_label: __address__
        replacement: <node hostname or IP>:9500
```

Then restart Prometheus and you're good to go!

## How does it work?

It calls the Coingecko API and returns it in the format Prometheus can consume.

## How can I configure it?

You can pass the artuments to the executable file to configure it. Here is the parameters list:

- `--listen_address` : Address on which to expose metrics and web interface, default `:9500`
- `--currency` : Coins currency, supported: USD, IDR, BTC, ETH, default `USD` 
- `--coingecko_api` : Coingecko API URL with its version, default `https://api.coingecko.com/api/v3`
- `--json` - output logs as JSON. Useful if you don't read it on servers but instead use logging aggregation solutions such as ELK stack.

## Which coins this is guaranteed to work?

It should work if coins already added in Coingecko, you can check it out Coingecko Coin Id in the detail page like this [one](https://www.coingecko.com/en/coins/cosmos-hub).

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.
