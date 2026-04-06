package nicru

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/n0isy/caddy-dns-nicru/nicrudns"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *nicrudns.Provider }

func init() {
	caddy.RegisterModule(&Provider{})
}

// CaddyModule returns the Caddy module information.
func (*Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.nicru",
		New: func() caddy.Module { return &Provider{new(nicrudns.Provider)} },
	}
}

// Provision implements the caddy.Provisioner interface — replaces placeholders.
func (p *Provider) Provision(_ caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.OAuth2ClientID = repl.ReplaceAll(p.Provider.OAuth2ClientID, "")
	p.Provider.OAuth2SecretID = repl.ReplaceAll(p.Provider.OAuth2SecretID, "")
	p.Provider.Username = repl.ReplaceAll(p.Provider.Username, "")
	p.Provider.Password = repl.ReplaceAll(p.Provider.Password, "")
	p.Provider.DnsServiceName = repl.ReplaceAll(p.Provider.DnsServiceName, "")
	p.Provider.CachePath = repl.ReplaceAll(p.Provider.CachePath, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	nicru {
//	    oauth2_client_id <id>
//	    oauth2_secret_id <secret>
//	    username <username>
//	    password <password>
//	    dns_service_name <service>
//	    cache_path <path>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "oauth2_client_id":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.OAuth2ClientID = d.Val()
			case "oauth2_secret_id":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.OAuth2SecretID = d.Val()
			case "username":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.Username = d.Val()
			case "password":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.Password = d.Val()
			case "dns_service_name":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.DnsServiceName = d.Val()
			case "cache_path":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.CachePath = d.Val()
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.OAuth2ClientID == "" {
		return d.Err("missing oauth2_client_id")
	}
	if p.Provider.OAuth2SecretID == "" {
		return d.Err("missing oauth2_secret_id")
	}
	if p.Provider.Username == "" {
		return d.Err("missing username")
	}
	if p.Provider.Password == "" {
		return d.Err("missing password")
	}
	if p.Provider.DnsServiceName == "" {
		return d.Err("missing dns_service_name")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
