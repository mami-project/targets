# MAMI Public Targets List

The H2020 Measurement and Architecture for a Middleboxed Internet (MAMI)
Project Public Targets List (PTL) is a set of Internet hostnames, resolved
addresses, and information about the protocol(s) these hosts are known to
listen on, to provide a set of targets for active Internet measurements. The
PTL is derived from the published results of ongoing and prior active Internet
measurement studies, as well as freely available lists of public measurement
targets. MAMI uses this list to provide destination diversity for Internet
path transparency studies, but we intend this resource to be of general
utility to all active measurement studies of services in the Internet.

## Targets

The target lists are provided as CSV files without header, with the following columns:

- 0: IP address (either IPv4 or IPv6)
- 1: Service number (80 = web, 25 = mail (listed as MX), 53 = dns (listed as NS))
- 2: An FQDN associated with the IP address

The most recent target lists were compiled and resolved on 7 December 2016:

- [Web servers](https://github.com/mami-project/targets/blob/master/public-targets-20161207-http.csv?raw=true)
- [Mail servers](https://github.com/mami-project/targets/blob/master/public-targets-20161207-smtp.csv?raw=true)
- [DNS servers](https://github.com/mami-project/targets/blob/master/public-targets-20161207-dns.csv?raw=true)

## Data Sources

We take targets from a variety of sources:

### Currently Integrated

- MAMI project ECN measurements, June 2016 (see [blog post](https://mami-project.eu/index.php/2016/06/13/70-of-popular-web-sites-support-ecn/) / [notebooks](https://github.com/mami-project/ecn-conspiracy) ) and December 2016 (results pending).

### In raw_sources, Not Yet Integrated

- [isthewebhttp2yet.com](https://isthewebhttp2yet.com), 16 November 2016
- [RIPE Atlas Anchors](https://atlas.ripe.net/about/anchors/), retrieved 22 November 2016

### Planned

- [Cisco Umbrella One Million](https://blog.opendns.com/2016/12/14/cisco-umbrella-1-million/); since this has a different source (DNS queries as opposed to tests that have already been run), will require different (probably minor) preprocessing.

