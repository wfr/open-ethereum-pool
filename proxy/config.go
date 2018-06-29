package proxy

import (
	"github.com/blockmaintain/open-ethereum-pool-nh/api"
	"github.com/blockmaintain/open-ethereum-pool-nh/payouts"
	"github.com/blockmaintain/open-ethereum-pool-nh/policy"
	"github.com/blockmaintain/open-ethereum-pool-nh/storage"
)

type Config struct {
	Name                  string        `json:"name"`
	Proxy                 Proxy         `json:"proxy"`
	Api                   api.ApiConfig `json:"api"`
	Upstream              []Upstream    `json:"upstream"`
	UpstreamCheckInterval string        `json:"upstreamCheckInterval"`

	Threads int `json:"threads"`

	Coin          string                 `json:"coin"`
	Redis         storage.Config         `json:"redis"`
	SQL           storage.SqlConfig      `json:"sql"`
	BlockUnlocker payouts.UnlockerConfig `json:"unlocker"`
	Payouts       payouts.PayoutsConfig  `json:"payouts"`

	NewrelicName    string `json:"newrelicName"`
	NewrelicKey     string `json:"newrelicKey"`
	NewrelicVerbose bool   `json:"newrelicVerbose"`
	NewrelicEnabled bool   `json:"newrelicEnabled"`
}

type Proxy struct {
	Enabled              bool    `json:"enabled"`
	Listen               string  `json:"listen"`
	LimitHeadersSize     int     `json:"limitHeadersSize"`
	LimitBodySize        int64   `json:"limitBodySize"`
	BehindReverseProxy   bool    `json:"behindReverseProxy"`
	BlockRefreshInterval string  `json:"blockRefreshInterval"`
	Difficulty           int64   `json:"difficulty"`
	DifficultyNiceHash   float64 `json:"difficultyNiceHash"`
	StateUpdateInterval  string  `json:"stateUpdateInterval"`
	HashrateExpiration   string  `json:"hashrateExpiration"`

	Policy policy.Config `json:"policy"`

	MaxFails    int64 `json:"maxFails"`
	HealthCheck bool  `json:"healthCheck"`

	Stratum Stratum `json:"stratum"`

	StratumNiceHash StratumNiceHash `json:"stratum_nice_hash"`
}

type Stratum struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout string `json:"timeout"`
	MaxConn int    `json:"maxConn"`
}

type StratumNiceHash struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout string `json:"timeout"`
	MaxConn int    `json:"maxConn"`
}

type Upstream struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Timeout string `json:"timeout"`
}
