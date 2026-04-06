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
                "oauth2_secret_id": "YOUR_SECRET_ID",
                "username": "YOUR_USERNAME",
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
            oauth2_secret_id {env.NICRU_SECRET_ID}
            username         {env.NICRU_USERNAME}
            password         {env.NICRU_PASSWORD}
            dns_service_name {env.NICRU_SERVICE_NAME}
            cache_path       /var/lib/caddy/.nicru-cache
        }
        propagation_delay 60s
    }
}
```

The `cache_path` directive is optional and specifies where to store the OAuth2 token cache. If omitted, tokens are re-fetched on each request.

Setting `propagation_delay` to 60s is recommended for NIC.RU as their DNS propagation can be slow.

## Authentication

To obtain API credentials:

1. Log in to [NIC.RU](https://www.nic.ru/).
2. Go to **Applications** in your account settings and register an OAuth2 application. Note the **Client ID** and **Secret**.
3. The `username` and `password` are your NIC.RU account credentials.
4. The `dns_service_name` is the name of your DNS hosting service (visible in the DNS management panel, e.g. `MY-DNS-SERVICE`).

Refer to the [NIC.RU DNS API documentation](https://www.nic.ru/help/upload/file/API_DNS-hosting.pdf) for details.
