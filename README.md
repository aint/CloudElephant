[![Build Status](https://github.com/aint/CloudElephant/actions/workflows/go.yml/badge.svg "GitHub Actions build status")](https://github.com/aint/CloudElephant/actions?query=workflow%3AGo)

<p align="center">
    <a href="https://www.youtube.com/watch?v=FoTYV22qZTg"><img src="https://i.imgur.com/G01TSPA.png" alt="Cloud Elephant" width="200"></a>
</p>

_Dedicated to Terry A. Davis. The smartest programmer that's ever lived._

# Cloud Elephant

Cloud Elephant is a tool providing a simple CLI interface for finding idle and unused resources in public clouds (AWS, Azure).

Supports:
 - [AWS ELB (Elastic Load Balancer)](#aws-elb)
 - [AWS EIP (Elastic IP Addresses)](#aws-eip)
 - [AWS EBS (Elastic Block Store)](#aws-ebs)
 - [AWS AMI (Machine Images)](#aws-ami)
 - AWS RDS (Relational Database Service) _planned_
 - AWS EC2 (Elastic Compute Cloud) _planned_
 - [Azure Load Balancer](#azure-load-balancer)
 - Azure Managed Disk _planned_

## Is it any good?
Yes.

## Installation

### MacOS

Use homebrew:

```
$ brew tap aint/cloudelephant-tap
$ brew install cloudelephant
```

### Build from sources

If you want the latest version, the recommended installation option is to use `go get`:

`$ go get -u github.com/aint/CloudElephant`

and add an alias:

`$ alias ce=CloudElephant`

### Download from the Releases page

Download a binary from the [GitHub Releases](https://github.com/aint/CloudElephant/releases) tab.


## Configuration

In order to use Azure, you need to set the following environment variables:

- `AZURE_SUBSCRIPTION_ID`
- `AZURE_TENANT_ID`
- `AZURE_CLIENT_ID`
- `AZURE_CLIENT_SECRET`


## Usage

`$ ce [unused|idle] [elb|elbv2|eip|ami|ebs|azlb]`

### AWS ELB

Find classic ELB with no associated back-end instances.

`$ ce unused elb`

Find ELBv2 (Application, Network, Gateway) which associated target groups has no EC2 target instance registered.

`$ ce unused elbv2`

### AWS EBS

Find available (unattached) EBS and EBS that are attached to stopped EC2 instances.

`$ ce unused ebs`

### AWS AMI

Find unused Amazon Machine Images (no instances are running from AMI).

`$ ce unused ami`

### AWS EIP

Find Elastic IP Addresses that is not associated with a running EC2 instance or an Elastic Network Interface.

`$ ce unused eip`

### Azure Load Balancer

Find Load Balancers which don't have any associated backend pool instances.

`$ ce unused azlb`
