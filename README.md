# Arbitrage

## Description

This is a simple implementation of arbitrage project where currently triangular arbitrage is implemented.  The Bellman-Ford algorithm is used to detect tradeable cycles. Right now, this is not a complete trading bot — there’s no order execution or risk management yet — but it’s a clean starting point for anyone who wants to step foot into implementing such project with Go language.

### Problem

Basically we are going to solve the following problem.

If there are three pairs and exchange rates like this:

* $r_{12}$: from currency USDT → BTC  
* $r_{23}$: from currency BTC → ETH  
* $r_{31}$: from currency ETH → USDT  

We also take fees into account:

$$
r^{\text{eff}}_{ij} = r_{ij} \cdot (1 - f_{ij})
$$

The condition for arbitrage is:

$$
r^{\text{eff}}_{12} \cdot r^{\text{eff}}_{23} \cdot r^{\text{eff}}_{31} > 1
$$

Now applying negative logarithm to both side:

$$
-(\log r^{\text{eff}}_{12} + \log r^{\text{eff}}_{23} + \log r^{\text{eff}}_{31}) < 0
$$

This condition corresponds to detecting negative cycles in the Bellman-Ford algorithm. In a Bellman-Ford graph with weights $w_{ij} = -\log(r_{ij})$, if the cycle $1 \to 2 \to 3 \to 1$ has a negative total weight, then there is an arbitrage opportunity.

## Limitations

* Only supports **Binance** exchange  
* Based on the REST API (slow compared to WebSockets)  
* No order execution / trading support  
* No risk management

## Features

* Uses [govalues/decimal](https://github.com/govalues/decimal) as a replacement for `float64`, with an added wrapper to support `Infinity`  
* Rate limiting via a custom Limiter to control HTTP request frequency to the exchange  
* Efficient use of Go goroutines to burst requests when needed (e.g., Depths, Klines, etc.)  
* Simple configuration file (TOML format)  

## Usage

Clone the repository:

```bash
git clone https://github.com/0xarash/arbitrage.git
cd arbitrage
```

Install dependencies:

```go
go mod tidy
```

Or build and run:

```go
go build -o arbitrage
./arbitrage
```

## Configuration

The project uses a simple config file (`config.toml`). Below is an example of such config file.

```toml
[arbitrage]
start_currencies = ["USDT"]

[ranking]
volume_threshold = 1_000_000
top_pairs = 500
ignore_zero = true

[worker]
concurrency = 8

[limiter]
weight_kline = 2
weight_depth = 5

[binance]
trading_fee = 0.001
```
