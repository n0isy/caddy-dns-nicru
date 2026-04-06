# NIC.RU API vs Client — Comparison & Roadmap

Validated against the official [NIC.RU DNS API docs](https://www.nic.ru/help/upload/file/API_DNS-hosting.pdf) and [OAuth server docs](https://www.nic.ru/help/oauth-server_3642.html).

## OAuth2 / Authentication

| API spec | Client | Status |
|---|---|---|
| `POST /oauth/token` | `TokenURL = BaseURL + "/oauth/token"` | Done |
| `grant_type=password` | `oauth2Config.PasswordCredentialsToken()` | Done |
| `client_id` + `client_secret` in params | `AuthStyle: oauth2.AuthStyleInParams` | Done |
| `scope=.+:/dns-master/.+` (full access) | `OAuth2Scope = ".+:/dns-master/.+"` | Done |
| `username=123/NIC-D` (contract ID) | Passed via `Provider.Username` | Done |
| `refresh_token` auto-refresh | `oauth2.Config.Client()` handles transparently | Done |

## Service Endpoints (Section 3)

| # | API endpoint | Method | Client function | Status |
|---|---|---|---|---|
| 3.1 | `/dns-master/services` | GET | `GetServices()` | Done |
| 3.2.1 | `/dns-master/zones` | GET | — | Not needed |
| 3.2.2 | `/dns-master/services/<svc>/zones` | GET | — | Not needed |
| 3.2.3 | `/dns-master/services/<svc>/zones/<zone>` | PUT | — | Roadmap |
| 3.2.4 | `/dns-master/zones/primary/<zone>` | PUT | — | Roadmap |
| 3.2.5 | `/dns-master/services/<svc>/zones/<zone>/move/<new>` | POST | — | Roadmap |
| 3.2.6 | `/dns-master/services/<svc>/zones/<zone>` | DELETE | — | Roadmap |
| 3.2.7 | `/dns-master/services/<svc>/zones/<zone>/xfer` | GET | — | Roadmap |
| 3.2.8 | `/dns-master/services/<svc>/zones/<zone>/xfer` | POST | — | Roadmap |

## DNS-master Record Operations (Section 4)

| # | API endpoint | Method | Client function | Status |
|---|---|---|---|---|
| 4.1.1 | `/dns-master/services/<svc>/zones/<zone>` | GET | `DownloadZone()` | Done |
| 4.1.2 | `/dns-master/services/<svc>/zones/<zone>` | POST | — | Roadmap |
| 4.2 | `/dns-master/services/<svc>/zones/<zone>/rollback` | POST | `RollbackZone()` | Done |
| 4.3.1 | `/dns-master/services/<svc>/zones/<zone>/commit` | POST | `CommitZone()` | Done |
| 4.3.2 | `/dns-master/services/<svc>/zones/<zone>/revisions` | GET | — | Roadmap |
| 4.3.3 | `/dns-master/services/<svc>/zones/<zone>/revisions/<rev>` | GET | — | Roadmap |
| 4.3.4 | `/dns-master/services/<svc>/zones/<zone>/revisions/<rev>` | PUT | — | Roadmap |
| 4.4.1 | `/dns-master/services/<svc>/zones/<zone>/default-ttl` | POST | — | Roadmap |
| 4.4.2 | `/dns-master/services/<svc>/zones/<zone>/default-ttl` | GET | — | Roadmap |
| 4.5.1 | `/dns-master/services/<svc>/zones/<zone>/records` | PUT | `Add()` | Done |
| 4.5.2 | `/dns-master/services/<svc>/zones/<zone>/records` | GET | `GetRecords()` | Done |
| 4.5.3 | `/dns-master/services/<svc>/zones/<zone>/records/<id>` | DELETE | `DeleteRecord()` | Done |

## XML Record Type Support (Section 4.6)

| # | RR Type | Read | Write | Status |
|---|---|---|---|---|
| 4.6.1 | SOA | model only | — | Read-only (auto-managed) |
| 4.6.2 | A | `GetARecords()` | `AddA()` | Done |
| 4.6.3 | AAAA | `GetAAAARecords()` | `AddAAAA()` | Done |
| 4.6.4 | CNAME | `GetCnameRecords()` | `AddCnames()` | Done |
| 4.6.5 | NS | model only | — | Roadmap |
| 4.6.6 | MX | `GetMxRecords()` | `AddMx()` | Done |
| 4.6.7 | SRV | model only | — | Roadmap |
| 4.6.8 | PTR | model only | — | Roadmap |
| 4.6.9 | TXT | `GetTxtRecords()` | `AddTxt()` | Done |
| 4.6.10 | DNAME | model only | — | Roadmap |
| 4.6.11 | HINFO | model only | — | Roadmap |
| 4.6.12 | NAPTR | model only | — | Roadmap |
| 4.6.13 | RP | model only | — | Roadmap |

## Secondary Service Endpoints (Section 5)

| # | API endpoint | Method | Client function | Status |
|---|---|---|---|---|
| 5.1.1 | `/dns-master/services/<svc>/zones/<zone>/masters` | GET | — | Roadmap |
| 5.1.2 | `/dns-master/services/<svc>/zones/<zone>/masters` | POST | — | Roadmap |

## ACME DNS-01 Challenge Flow

The minimum viable path for Let's Encrypt wildcard certificates:

| Step | API call | Client | Status |
|---|---|---|---|
| 1. Authenticate | `POST /oauth/token` | `GetOauth2Client()` | Done |
| 2. Add `_acme-challenge` TXT | `PUT .../records` | `AddTxt()` → `Add()` | Done |
| 3. Commit zone to DNS servers | `POST .../commit` | `CommitZone()` | Done |
| 4. Wait for propagation | — | Caddy `propagation_delay` | External |
| 5. Delete challenge TXT | `DELETE .../records/<id>` | `DeleteRecord()` | Done |
| 6. Commit again | `POST .../commit` | `CommitZone()` | Done |

## Changelog

### 2026-04-06 — Initial release

Bugs fixed from upstream [libdns/nicrudns](https://github.com/libdns/nicrudns):

| File | Bug | Impact |
|---|---|---|
| `add-txt.go` | `Type: "MX"` instead of `Type: "TXT"` | Critical — ACME DNS-01 would never work |
| `rollback.go` | Missing error check after `client.Do()` | Medium — nil pointer panic on network error |
| `cache.go` | No guard for empty `CachePath` | Medium — crash when `cache_path` not set |

Other changes:
- Rewrote `provider.go` for libdns v1.1 API (interface-based `Record` instead of struct)
- Separated A/AAAA into distinct `AddA()`/`AddAAAA()` calls
- Correct MX data parsing (`"preference target"` → separate fields)
