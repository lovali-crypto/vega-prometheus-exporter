package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Validators []struct {
	Address string `json:"address"`
	PubKey  struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"pub_key"`
	Power string `json:"power"`
	Name  string `json:"name"`
}

type VegaValidator struct {
	Name         string
	Address      string
	ShortAddress string
}

type VegaStatus struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		NodeInfo struct {
			ProtocolVersion struct {
				P2P   string `json:"p2p"`
				Block string `json:"block"`
				App   string `json:"app"`
			} `json:"protocol_version"`
			ID         string `json:"id"`
			ListenAddr string `json:"listen_addr"`
			Network    string `json:"network"`
			Version    string `json:"version"`
			Channels   string `json:"channels"`
			Moniker    string `json:"moniker"`
			Other      struct {
				TxIndex    string `json:"tx_index"`
				RPCAddress string `json:"rpc_address"`
			} `json:"other"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHash     string    `json:"latest_block_hash"`
			LatestAppHash       string    `json:"latest_app_hash"`
			LatestBlockHeight   string    `json:"latest_block_height"`
			LatestBlockTime     time.Time `json:"latest_block_time"`
			EarliestBlockHash   string    `json:"earliest_block_hash"`
			EarliestAppHash     string    `json:"earliest_app_hash"`
			EarliestBlockHeight string    `json:"earliest_block_height"`
			EarliestBlockTime   time.Time `json:"earliest_block_time"`
			CatchingUp          bool      `json:"catching_up"`
		} `json:"sync_info"`
		ValidatorInfo struct {
			Address string `json:"address"`
			PubKey  struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
			VotingPower string `json:"voting_power"`
		} `json:"validator_info"`
	} `json:"result"`
}

type VegaConsensus struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		RoundState struct {
			Height     string    `json:"height"`
			Round      int       `json:"round"`
			Step       int       `json:"step"`
			StartTime  time.Time `json:"start_time"`
			CommitTime time.Time `json:"commit_time"`
			Validators struct {
				Validators []struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"validators"`
				Proposer struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"proposer"`
			} `json:"validators"`
			Proposal           interface{} `json:"proposal"`
			ProposalBlock      interface{} `json:"proposal_block"`
			ProposalBlockParts interface{} `json:"proposal_block_parts"`
			LockedRound        int         `json:"locked_round"`
			LockedBlock        interface{} `json:"locked_block"`
			LockedBlockParts   interface{} `json:"locked_block_parts"`
			ValidRound         int         `json:"valid_round"`
			ValidBlock         interface{} `json:"valid_block"`
			ValidBlockParts    interface{} `json:"valid_block_parts"`
			Votes              []struct {
				Round              int      `json:"round"`
				Prevotes           []string `json:"prevotes"`
				PrevotesBitArray   string   `json:"prevotes_bit_array"`
				Precommits         []string `json:"precommits"`
				PrecommitsBitArray string   `json:"precommits_bit_array"`
			} `json:"votes"`
			CommitRound int `json:"commit_round"`
			LastCommit  struct {
				Votes         []interface{} `json:"votes"`
				VotesBitArray string        `json:"votes_bit_array"`
				PeerMaj23S    struct {
				} `json:"peer_maj_23s"`
			} `json:"last_commit"`
			LastValidators struct {
				Validators []struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"validators"`
				Proposer struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"proposer"`
			} `json:"last_validators"`
			TriggeredTimeoutPrecommit bool `json:"triggered_timeout_precommit"`
		} `json:"round_state"`
		Peers []struct {
			NodeAddress string `json:"node_address"`
			PeerState   struct {
				RoundState struct {
					Height                     string    `json:"height"`
					Round                      int       `json:"round"`
					Step                       int       `json:"step"`
					StartTime                  time.Time `json:"start_time"`
					Proposal                   bool      `json:"proposal"`
					ProposalBlockPartSetHeader struct {
						Total int    `json:"total"`
						Hash  string `json:"hash"`
					} `json:"proposal_block_part_set_header"`
					ProposalBlockParts interface{} `json:"proposal_block_parts"`
					ProposalPolRound   int         `json:"proposal_pol_round"`
					ProposalPol        string      `json:"proposal_pol"`
					Prevotes           string      `json:"prevotes"`
					Precommits         string      `json:"precommits"`
					LastCommitRound    int         `json:"last_commit_round"`
					LastCommit         string      `json:"last_commit"`
					CatchupCommitRound int         `json:"catchup_commit_round"`
					CatchupCommit      string      `json:"catchup_commit"`
				} `json:"round_state"`
				Stats struct {
					Votes      string `json:"votes"`
					BlockParts string `json:"block_parts"`
				} `json:"stats"`
			} `json:"peer_state"`
		} `json:"peers"`
	} `json:"result"`
}

type VegaNetInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Listening bool     `json:"listening"`
		Listeners []string `json:"listeners"`
		NPeers    string   `json:"n_peers"`
		Peers     []struct {
			NodeInfo struct {
				ProtocolVersion struct {
					P2P   string `json:"p2p"`
					Block string `json:"block"`
					App   string `json:"app"`
				} `json:"protocol_version"`
				ID         string `json:"id"`
				ListenAddr string `json:"listen_addr"`
				Network    string `json:"network"`
				Version    string `json:"version"`
				Channels   string `json:"channels"`
				Moniker    string `json:"moniker"`
				Other      struct {
					TxIndex    string `json:"tx_index"`
					RPCAddress string `json:"rpc_address"`
				} `json:"other"`
			} `json:"node_info"`
			IsOutbound       bool `json:"is_outbound"`
			ConnectionStatus struct {
				Duration    string `json:"Duration"`
				SendMonitor struct {
					Start    time.Time `json:"Start"`
					Bytes    string    `json:"Bytes"`
					Samples  string    `json:"Samples"`
					InstRate string    `json:"InstRate"`
					CurRate  string    `json:"CurRate"`
					AvgRate  string    `json:"AvgRate"`
					PeakRate string    `json:"PeakRate"`
					BytesRem string    `json:"BytesRem"`
					Duration string    `json:"Duration"`
					Idle     string    `json:"Idle"`
					TimeRem  string    `json:"TimeRem"`
					Progress int       `json:"Progress"`
					Active   bool      `json:"Active"`
				} `json:"SendMonitor"`
				RecvMonitor struct {
					Start    time.Time `json:"Start"`
					Bytes    string    `json:"Bytes"`
					Samples  string    `json:"Samples"`
					InstRate string    `json:"InstRate"`
					CurRate  string    `json:"CurRate"`
					AvgRate  string    `json:"AvgRate"`
					PeakRate string    `json:"PeakRate"`
					BytesRem string    `json:"BytesRem"`
					Duration string    `json:"Duration"`
					Idle     string    `json:"Idle"`
					TimeRem  string    `json:"TimeRem"`
					Progress int       `json:"Progress"`
					Active   bool      `json:"Active"`
				} `json:"RecvMonitor"`
				Channels []struct {
					ID                int    `json:"ID"`
					SendQueueCapacity string `json:"SendQueueCapacity"`
					SendQueueSize     string `json:"SendQueueSize"`
					Priority          string `json:"Priority"`
					RecentlySent      string `json:"RecentlySent"`
				} `json:"Channels"`
			} `json:"connection_status"`
			RemoteIP string `json:"remote_ip"`
		} `json:"peers"`
	} `json:"result"`
}

const namespace = "vega"
const vegaStatusUrl = "/status"
const vegaConsensusUrl = "/dump_consensus_state"
const vegaGenesisUrl = "/genesis"
const netInfo = "/net_info"

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}

	listenAddress = flag.String("web.listen-address", ":9141",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last vega query successful.",
		nil, nil,
	)
	metricCatchingUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "sync_cytching_up"),
		"Is the node catching up?",
		nil, nil,
	)
	metricValidatorSigning = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "validator_signing"),
		"Flag indicating if a validator is signing or not (per validator).",
		[]string{"validator"}, nil,
	)
)

type Exporter struct {
	vegaEndpoint string
}

func NewExporter(vegaEndpoint string) *Exporter {
	return &Exporter{
		vegaEndpoint: vegaEndpoint,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- metricCatchingUp
	ch <- metricValidatorSigning
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	_, err := e.LoadVegaStatus(ch)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	validators, err := e.GetVegaValidators()

	err = e.LoadVegaConsensus(validators, ch)
}

func (e *Exporter) LoadVegaStatus(ch chan<- prometheus.Metric) (VegaStatus, error) {
	// we initialize our array
	var vegaStatus VegaStatus
	req, err := http.NewRequest("GET", e.vegaEndpoint+vegaStatusUrl, nil)
	if err != nil {
		return vegaStatus, err
	}

	// Make request and show output.
	resp, err := client.Do(req)
	if err != nil {
		return vegaStatus, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return vegaStatus, err
	}
	//fmt.Println(string(body))

	// we unmarshal our byteArray which contains our
	// json content into 'vegaStatus' which we defined above
	err = json.Unmarshal(body, &vegaStatus)
	if err != nil {
		return vegaStatus, err
	}

	var catching float64
	catching = 0

	if vegaStatus.Result.SyncInfo.CatchingUp == true {
		catching = 1
	}

	ch <- prometheus.MustNewConstMetric(
		metricCatchingUp, prometheus.GaugeValue, catching,
	)

	return vegaStatus, nil
}

func (e *Exporter) GetVegaValidators() ([]VegaValidator, error) {
	// Get Vega genesis file
	req, err := http.NewRequest("GET", e.vegaEndpoint+netInfo, nil)
	if err != nil {
		return nil, err
	}

	// Make request and show output.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	var validators VegaNetInfo
	v, err := json.Marshal(result["result"])
	err = json.Unmarshal(v, &result)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(v, &validators)
	//log.Printf("result: %+v\n", result)
	//log.Printf("marshaled result: %+v\n", v)

	var retValidators []VegaValidator
	for _, val := range validators.Result.Peers {
		var validator VegaValidator
		validator.Name = val.NodeInfo.Moniker
		validator.Address = val.NodeInfo.ID
		validator.ShortAddress = val.NodeInfo.ID[0:12]
		retValidators = append(retValidators, validator)
	}

	//log.Printf("validators: %+v\n", validators)

	return retValidators, nil
}

func (e *Exporter) LoadVegaConsensus(validators []VegaValidator, ch chan<- prometheus.Metric) error {
	var vegaConsensus VegaConsensus
	// Load channel stats
	req, err := http.NewRequest("GET", e.vegaEndpoint+vegaConsensusUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make request and show output.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(body))
	err = json.Unmarshal(body, &vegaConsensus)
	if err != nil {
		return err
	}

	votes := GetVoteSlice(vegaConsensus.Result.RoundState.LastCommit.Votes)
	log.Printf("%+v\n", votes)
	log.Printf("%+v\n", validators)

	for _, val := range validators {
		//log.Printf("Parsing validator %s\n", val.Name)
		if contains(votes, val.ShortAddress) {
			ch <- prometheus.MustNewConstMetric(
				metricValidatorSigning, prometheus.GaugeValue, 1, val.Name,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				metricValidatorSigning, prometheus.GaugeValue, 0, val.Name,
			)
		}
	}

	log.Println("Endpoint scraped")
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		log.Printf("'%s' '%s'\n", a, e)
		if strings.TrimSpace(a) == strings.TrimSpace(e) {
			return true
		}
	}
	return false
}

func GetVoteSlice(votesInt []interface{}) []string {
	var votes []string
	for _, val := range votesInt {
		str := fmt.Sprintf("%v", val)
		re := regexp.MustCompile("([0-9A-Z])+ ")
		match := re.FindStringSubmatch(str)
		if match != nil {
			//fmt.Println(match[0])
			votes = append(votes, match[0])
		}
	}
	log.Println(votes)
	return votes
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assume env variables are set.")
	}

	flag.Parse()

	vegaEndpoint := os.Getenv("VEGA_ENDPOINT")

	exporter := NewExporter(vegaEndpoint)
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Vega Metrics Exporter</title></head>
             <body>
             <h1>Vega Metrics Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
