package nicrudns

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/pkg/errors"
)

// Provider facilitates DNS record manipulation with NIC.ru.
type Provider struct {
	OAuth2ClientID string `json:"oauth2_client_id"`
	OAuth2SecretID string `json:"oauth2_secret_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DnsServiceName string `json:"dns_service_name"`
	CachePath      string `json:"cache_path"`
}

func parseTTL(ttlStr string) time.Duration {
	if v, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
		return time.Duration(v) * time.Second
	}
	return 0
}

// parseMXData splits MX RR data "preference target" into separate fields.
func parseMXData(data string) (preference, target string) {
	parts := strings.Fields(data)
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "10", data
}

// rrToRecord converts an internal RR to a libdns Record (v1.1 interface).
func rrToRecord(rr *RR) libdns.Record {
	ttl := parseTTL(rr.Ttl)
	switch {
	case rr.A != nil:
		ip, _ := netip.ParseAddr(rr.A.String())
		return libdns.Address{Name: rr.Name, TTL: ttl, IP: ip}
	case rr.AAAA != nil:
		ip, _ := netip.ParseAddr(rr.AAAA.String())
		return libdns.Address{Name: rr.Name, TTL: ttl, IP: ip}
	case rr.Cname != nil:
		return libdns.CNAME{Name: rr.Name, TTL: ttl, Target: rr.Cname.Name}
	case rr.Txt != nil:
		return libdns.TXT{Name: rr.Name, TTL: ttl, Text: rr.Txt.String}
	case rr.Mx != nil:
		pref, _ := strconv.ParseUint(rr.Mx.Preference, 10, 16)
		target := ""
		if rr.Mx.Exchange != nil {
			target = rr.Mx.Exchange.Name
		}
		return libdns.MX{Name: rr.Name, TTL: ttl, Preference: uint16(pref), Target: target}
	default:
		return libdns.RR{Name: rr.Name, TTL: ttl, Type: rr.Type}
	}
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	client := NewClient(p)
	rrs, err := client.GetRecords(zone)
	if err != nil {
		return nil, err
	}
	var records []libdns.Record
	for _, rr := range rrs {
		records = append(records, rrToRecord(rr))
	}
	return records, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client := NewClient(p)
	var result []libdns.Record
	for _, rec := range records {
		rr := rec.RR()
		ttlStr := strconv.Itoa(int(rr.TTL.Seconds()))
		switch rr.Type {
		case "A":
			response, err := client.AddA(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "AAAA":
			response, err := client.AddAAAA(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "CNAME":
			response, err := client.AddCnames(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "TXT":
			response, err := client.AddTxt(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "MX":
			pref, target := parseMXData(rr.Data)
			response, err := client.AddMx(zone, []string{rr.Name}, target, pref, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		default:
			return nil, errors.Wrap(NotImplementedRecordType, rr.Type)
		}
	}
	if _, err := client.CommitZone(zone); err != nil {
		return nil, err
	}
	return result, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client := NewClient(p)
	allRecords, err := client.GetRecords(zone)
	if err != nil {
		return nil, err
	}
	var result []libdns.Record
	for _, rec := range records {
		rr := rec.RR()
		ttlStr := strconv.Itoa(int(rr.TTL.Seconds()))

		// Delete existing records with matching name+type
		for _, existing := range allRecords {
			if existing.Name == rr.Name && existing.Type == rr.Type {
				id, err := strconv.ParseInt(existing.ID, 10, 64)
				if err != nil {
					return nil, err
				}
				if _, err := client.DeleteRecord(zone, int(id)); err != nil {
					return nil, err
				}
			}
		}

		switch rr.Type {
		case "A":
			response, err := client.AddA(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "AAAA":
			response, err := client.AddAAAA(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "CNAME":
			response, err := client.AddCnames(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "TXT":
			response, err := client.AddTxt(zone, []string{rr.Name}, rr.Data, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		case "MX":
			pref, target := parseMXData(rr.Data)
			response, err := client.AddMx(zone, []string{rr.Name}, target, pref, ttlStr)
			if err != nil {
				return nil, err
			}
			result = append(result, rrToRecord(response.Data.Zone[0].Rr[0]))
		default:
			return nil, errors.Wrap(NotImplementedRecordType, rr.Type)
		}
	}
	if _, err := client.CommitZone(zone); err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteRecords deletes the records from the zone.
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	client := NewClient(p)
	allRecords, err := client.GetRecords(zone)
	if err != nil {
		return nil, err
	}
	var result []libdns.Record
	for _, rec := range records {
		rr := rec.RR()
		for _, existing := range allRecords {
			if existing.Name == rr.Name && existing.Type == rr.Type {
				id, err := strconv.ParseInt(existing.ID, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid record ID %q: %w", existing.ID, err)
				}
				if _, err := client.DeleteRecord(zone, int(id)); err != nil {
					return nil, err
				}
				result = append(result, rec)
				break
			}
		}
	}
	if _, err := client.CommitZone(zone); err != nil {
		return nil, err
	}
	return result, nil
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
