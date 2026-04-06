# NIC.RU (RU-CENTER) Module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It manages DNS records for domains hosted at [NIC.RU (RU-CENTER)](https://www.nic.ru/) using their [DNS API](https://www.nic.ru/help/upload/file/API_DNS-hosting.pdf). Based on the [libdns/nicrudns](https://github.com/libdns/nicrudns) provider (vendored and updated for libdns v1.1).

## Caddy Module Name

```
dns.providers.nicru
```

## Building

```bash
xcaddy build --with github.com/n0isy/caddy-dns-nicru@main
```

## Configuration

To use this module for the ACME DNS challenge, configure the [ACME issuer](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) in your Caddy JSON as follows:

```json
{
    "module": "acme",
    "challenges": {
        "dns": {
            "provider": {
                "name": "nicru",
                "oauth2_client_id": "YOUR_CLIENT_ID",
                "oauth2_secret_id": "YOUR_CLIENT_SECRET",
                "username": "123/NIC-D",
                "password": "YOUR_PASSWORD",
                "dns_service_name": "YOUR_SERVICE_NAME"
            }
        }
    }
}
```

Or in the Caddyfile:

```
your.domain.ru {
    respond "Hello World"

    tls {
        dns nicru {
            oauth2_client_id {env.NICRU_CLIENT_ID}
            oauth2_secret_id {env.NICRU_CLIENT_SECRET}
            username         {env.NICRU_USERNAME}
            password         {env.NICRU_PASSWORD}
            dns_service_name {env.NICRU_SERVICE_NAME}
            cache_path       /var/lib/caddy/.nicru-cache
        }
        propagation_delay 60s
    }
}
```

The `cache_path` directive is optional and specifies where to store the OAuth2 token cache file. If omitted, tokens are re-fetched on each request.

Setting `propagation_delay` to 60s is recommended for NIC.RU as their DNS propagation can be slow.

## Authentication

NIC.RU uses [OAuth 2.0](https://www.nic.ru/help/oauth-server_3642.html) with the Resource Owner Password Credentials grant (`grant_type=password`).

### 1. Register an OAuth2 application

1. Log in to [NIC.RU](https://www.nic.ru/).
2. Go to the [application registration page](https://www.nic.ru/manager/oauth.cgi?step=oauth.app_register).
3. Enter your application name and click **Register**.
4. The server will generate `client_id` and `client_secret` — these map to `oauth2_client_id` and `oauth2_secret_id` in the Caddy config.

You can manage your apps and rotate `client_secret` at [application management](https://www.nic.ru/manager/oauth.cgi?step=oauth.app_list).

### 2. Provide your NIC.RU credentials

- `username` — your NIC.RU contract identifier (e.g. `123/NIC-D` or `456/NIC-REG`). This is your login, not your email.
- `password` — your NIC.RU account password (administrative or technical password).

### 3. Find your DNS service name

- `dns_service_name` — the name of your DNS hosting service, visible in the [DNS management panel](https://www.nic.ru/manager/dns/serviceslist.cgi) (e.g. `MY-DNS-SERVICE`).

### Token flow

Under the hood, the provider obtains an access token via:

```
POST https://api.nic.ru/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password
&username=123/NIC-D
&password=<password>
&client_id=<client_id>
&client_secret=<client_secret>
&scope=.*
```

The token is cached (if `cache_path` is set) and automatically refreshed using the `refresh_token` when it expires.

## References

- [NIC.RU OAuth server docs](https://www.nic.ru/help/oauth-server_3642.html)
- [NIC.RU DNS API docs](https://www.nic.ru/help/upload/file/API_DNS-hosting.pdf)
- [libdns/nicrudns](https://github.com/libdns/nicrudns)
