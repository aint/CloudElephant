[![Build Status](https://github.com/aint/CloudElephant/workflows/Go/badge.svg "GitHub Actions build status")](https://github.com/aint/CloudElephant/actions?query=workflow%3AGo)

<p align="center">
    <a href="https://www.youtube.com/watch?v=FoTYV22qZTg"><img src="https://i.imgur.com/G01TSPA.png" alt="Cloud Elephant" width="200"></a>
</p>

_Dedicated to Terry A. Davis. The smartest programmer that's ever lived._

# Cloud Elephant

Cloud Elephant is a tool providing a simple CLI interface for finding idle and unused resources in public clouds (AWS, Azure).

Supports:
 - AWS ELB (Elastic Load Balancer)
 - AWS EIP (Elastic IP Addresses)
 - AWS EBS (Elastic Block Store)
 - AWS AMI (Machine Images)
 - AWS RDS (Relational Database Service) _planned_
 - AWS EC2 (Elastic Compute Cloud) _planned_
 - Azure Managed Disk _planned_
 - Azure Load Balancer _planned_

## Installation

### Build from sources

Use `go get` to get the latest version:

`$ go get -u github.com/aint/CloudElephant`

and add an alias:

`$ alias ce=CloudElephant`

### Download from the Releases page

Soon.

## Is it any good?
Yes.
