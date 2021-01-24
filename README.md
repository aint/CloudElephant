[![Build Status](https://github.com/aint/CloudElephant/workflows/Go/badge.svg "GitHub Actions build status")](https://github.com/aint/CloudElephant/actions?query=workflow%3AGo)

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
 - Azure Managed Disk _planned_
 - Azure Load Balancer _planned_

## Installation

### Build from sources

If you want the latest version, the recommended installation option is to use `go get`:

`$ go get -u github.com/aint/CloudElephant`

and add an alias:

`$ alias ce=CloudElephant`

### Download from the Releases page

Download a binary from the [GitHub Releases](https://github.com/aint/CloudElephant/releases) tab.

## Is it any good?
Yes.

## Usage

`$ ce [unused|idle] [elb|eip|ami|ebs|ec2|rds|az_disk|az_lb]`

### AWS ELB

Find ELB with no associated back-end instances and ELBv2 which associated target groups has no EC2 target instance registered.

`$ ce idle elb`

### AWS EBS

Find available (unattached) EBS and EBS that are attached to stopped EC2 instances.

`$ ce idle ebs`

### AWS AMI

Find unused Amazon Machine Images (no instances are running from AMI).

`$ ce idle ami`

### AWS EIP

Find Elastic IP Addresses that is not associated with a running EC2 instance or an Elastic Network Interface.

`$ ce idle eip`
