package pkg

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/smartwalle/nconv"
	"github.com/smartwalle/net4go"
	"github.com/smartwalle/xid"
)

type ServerConfig struct {
	Name      string `ini:"name"         json:"name"          yaml:"name"`
	LogStdout bool   `ini:"log_stdout"   json:"log_stdout"    yaml:"log_stdout"`
	LogFile   bool   `ini:"log_file"     json:"log_file"      yaml:"log_file"`
}

type ETCDConfig struct {
	Endpoints []string `ini:"endpoints"  json:"endpoints"    yaml:"endpoints"`
	XIDNode   string   `ini:"xid_node"   json:"xid_node"     yaml:"xid_node"`
}

type SQLConfig struct {
	Driver  string `ini:"driver"         json:"driver"        yaml:"driver"`
	URL     string `ini:"url"            json:"url"           yaml:"url"`
	MaxOpen int    `ini:"max_open"       json:"max_open"      yaml:"max_open"`
	MaxIdle int    `ini:"max_idle"       json:"max_idle"      yaml:"max_idle"`
}

type MongoConfig struct {
	URL string `ini:"url"     json:"url"      yaml:"url"`
}

type RedisConfig struct {
	Addr         string `ini:"addr"             json:"addr"             yaml:"addr"`
	Password     string `ini:"password"         json:"password"         yaml:"password"`
	DB           int    `ini:"db"               json:"db"               yaml:"db"`
	PoolSize     int    `ini:"pool_size"        json:"pool_size"        yaml:"pool_size"`
	MinIdleConns int    `ini:"min_idle_conns"   json:"min_idle_conns"   yaml:"min_idle_conns"`
}

type GRPCConfig struct {
	Domain        string `ini:"domain"             json:"domain"              yaml:"domain"`
	Name          string `ini:"name"               json:"name"                yaml:"name"`
	Node          string `ini:"node"               json:"node"                yaml:"node"`
	IP            string `ini:"ip"                 json:"ip"                  yaml:"ip"`
	Port          string `ini:"port"               json:"port"                yaml:"port"`
	ClientTracing bool   `ini:"client_tracing"     json:"client_tracing"      yaml:"client_tracing"`
	ServerTracing bool   `ini:"server_tracing"     json:"server_tracing"      yaml:"server_tracing"`
	GracefulStop  bool   `ini:"graceful_stop"      json:"graceful_stop"       yaml:"graceful_stop"`
}

func (cfg *GRPCConfig) GetAddress() string {
	if cfg.IP == "" {
		var err error
		cfg.IP, err = net4go.GetInternalIP()
		if err != nil {
			slog.Error("获取本地 IP 地址发生错误", slog.Any("error", err))
			os.Exit(-1)
		}
	}
	return fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)
}

func (cfg *GRPCConfig) GetDomain() string {
	return cfg.Domain
}

func (cfg *GRPCConfig) GetName() string {
	return cfg.Name
}

func (cfg *GRPCConfig) GetNode() string {
	if cfg.Node == "" {
		cfg.Node = fmt.Sprintf("%d", xid.Next())
	} else {
		cfg.Node = fmt.Sprintf("%s/%d", cfg.Node, xid.Next())
	}
	return cfg.Node
}

func (cfg *GRPCConfig) GetService() string {
	return fmt.Sprintf("%s/%s", cfg.Domain, cfg.Name)
}

type HTTPConfig struct {
	IP          string `ini:"ip"                json:"ip"                 yaml:"ip"`
	Port        string `ini:"port"              json:"port"               yaml:"port"`
	Name        string `ini:"name"              json:"name"               yaml:"name"`
	Domain      string `ini:"domain"            json:"domain"             yaml:"domain"`
	PPROFPath   string `ini:"pprof_path"        json:"pprof_path"         yaml:"pprof_path"`
	SwaggerPath string `ini:"swagger_path"      json:"swagger_path"       yaml:"swagger_path"`
}

func (cfg *HTTPConfig) Address() string {
	if cfg.IP == "" {
		var err error
		cfg.IP, err = net4go.GetInternalIP()
		if err != nil {
			slog.Error("获取本地 IP 地址发生错误", slog.Any("error", err))
			os.Exit(-1)
		}
	}
	if cfg.IP == "-" {
		cfg.IP = ""
	}
	if cfg.Port == "" || cfg.Port == "0" {
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", cfg.IP))
		if err != nil {
			slog.Error("获取随机端口发生错误", slog.Any("error", err))
			os.Exit(-1)
		}
		cfg.Port = nconv.String(listener.Addr().(*net.TCPAddr).Port)
		listener.Close()
	}
	return fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)
}

type TCPConfig struct {
	IP              string `ini:"ip"                        json:"ip"                   yaml:"ip"`
	Port            string `ini:"port"                      json:"port"                 yaml:"port"`
	MaxPacketSize   int    `ini:"max_packet_size"           json:"max_packet_size"      yaml:"max_packet_size"`
	ReadTimeout     int    `ini:"read_timeout"              json:"read_timeout"         yaml:"read_timeout"`
	WriteTimeout    int    `ini:"write_timeout"             json:"write_timeout"        yaml:"write_timeout"`
	ConnPerSecond   int    `ini:"conn_per_second"           json:"conn_per_second"      yaml:"conn_per_second"`
	PacketPerSecond int    `ini:"packet_per_second"         json:"packet_per_second"    yaml:"packet_per_second"`
}

func (cfg *TCPConfig) Address() string {
	if cfg.IP == "" {
		var err error
		cfg.IP, err = net4go.GetInternalIP()
		if err != nil {
			slog.Error("获取本地 IP 地址发生错误", slog.Any("error", err))
			os.Exit(-1)
		}
	}
	return fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)
}

type WebSocketConfig struct {
	IP              string `ini:"ip"                       json:"ip"                   yaml:"ip"`
	Port            string `ini:"port"                     json:"port"                 yaml:"port"`
	ReadTimeout     int    `ini:"read_timeout"             json:"read_timeout"         yaml:"read_timeout"`
	WriteTimeout    int    `ini:"write_timeout"            json:"write_timeout"        yaml:"write_timeout"`
	ReadBufferSize  int    `ini:"read_buffer_size"         json:"read_buffer_size"     yaml:"read_buffer_size"`
	WriteBufferSize int    `ini:"write_buffer_size"        json:"write_buffer_size"    yaml:"write_buffer_size"`
}

func (cfg *WebSocketConfig) Address() string {
	if cfg.IP == "" {
		var err error
		cfg.IP, err = net4go.GetInternalIP()
		if err != nil {
			slog.Error("获取本地 IP 地址发生错误", slog.Any("error", err))
			os.Exit(-1)
		}
	}
	return fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)
}

type NSQConfig struct {
	NSQAddr        string   `ini:"nsq_addr"                json:"nsq_addr"              yaml:"nsq_addr"`
	NSQLookupAddrs []string `ini:"nsq_lookup_addrs"        json:"nsq_lookup_addrs"      yaml:"nsq_lookup_addrs"`
	Group          string   `ini:"group"                   json:"group"                 yaml:"group"`
}

type ApplePayConfig struct {
	BundleId   string `ini:"bundle_id"                json:"bundle_id"               yaml:"bundle_id"`
	Sandbox    bool   `ini:"sandbox"                  json:"sandbox"                 yaml:"sandbox"`
	Production bool   `ini:"production"               json:"production"              yaml:"production"`
}

type AliPayConfig struct {
	AppId            string `ini:"app_id"                      json:"app_id"                  yaml:"app_id"`
	PrivateKey       string `ini:"private_key"                 json:"private_key"             yaml:"private_key"`
	IsProduction     bool   `ini:"is_production"               json:"is_production"           yaml:"is_production"`
	AppPublicCert    string `ini:"app_public_cert"             json:"app_public_cert"         yaml:"app_public_cert"`
	AliPayRootCert   string `ini:"ali_pay_root_cert"           json:"ali_pay_root_cert"       yaml:"ali_pay_root_cert"`
	AliPayPublicCert string `ini:"ali_pay_public_cert"         json:"ali_pay_public_cert"     yaml:"ali_pay_public_cert"`
	NotifyURL        string `ini:"notify_url"                  json:"notify_url"              yaml:"notify_url"`
	ReturnURL        string `ini:"return_url"                  json:"return_url"              yaml:"return_url"`
}

type WXPayConfig struct {
	AppId          string `ini:"app_id"                     json:"app_id"                   yaml:"app_id"`
	MchId          string `ini:"mch_id"                     json:"mch_id"                   yaml:"mch_id"`
	MchCertSN      string `ini:"mch_cert_sn"                json:"mch_cert_sn"              yaml:"mch_cert_sn"`
	MchAPIKey      string `ini:"mch_api_key"                json:"mch_api_key"              yaml:"mch_api_key"`
	WXPayClientKey string `ini:"wx_pay_client_key"          json:"wx_pay_client_key"        yaml:"wx_pay_client_key"`
	NotifyURL      string `ini:"notify_url"                 json:"notify_url"               yaml:"notify_url"`
}

type AliOSSConfig struct {
	Endpoint string `ini:"endpoint"            json:"endpoint"         yaml:"endpoint"`
	Key      string `ini:"key"                 json:"key"              yaml:"key"`
	Secret   string `ini:"secret"              json:"secret"           yaml:"secret"`
	Bucket   string `ini:"bucket"              json:"bucket"           yaml:"bucket"`
	Domain   string `ini:"domain"              json:"domain"           yaml:"domain"`
}
