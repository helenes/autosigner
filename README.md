# Autosigner

This app determines if a Certificate Request should be signed by the Puppet Server.

## Required configuration

### Autosigner config

Autosigner uses a configuration file [autosigner.yaml](configs/autosigner.yaml) for defining what GCP projects and AWS accounts are allowed certificate signing.

#### Location

Autosigner will look for the configuration file in two places: `/etc/puppetlabs/puppet/` and the same directory as the binary.
This configuration file is not automatically created by the Autosigner app.

#### Format

The [autosigner.yaml](configs/autosigner.yaml) file uses the [YAML](https://yaml.org/) format.

### IAM Permissions and Roles

> Note: TODO

## Methods for Signing

There are 3 methods for determining if a certificate should be signed.

- [x] Basic Autosigning
- [x] Google GCE
- [ ] Amazon EC2

### Basic Autosigning

Similar to [autosign.conf](https://puppet.com/docs/puppet/latest/ssl_autosign.html#ssl_basic_autosigning), Autosigner will use a file
called [autosigner_hostnames.conf](configs/autosigner_hostnames.conf)

#### Location

Autosigner will look for the file in the same path as the app. So if the Autosigner binary is in `/etc/puppetlabs/puppet/autosigner`, then
the location of the file should be `/etc/puppetlabs/puppet/autosigner_hostnames.conf`. This file is not automatically created by the Autosigner
app.

#### Format

The [autosigner_hostnames.conf](configs/autosigner_hostnames.conf) file is a line-separated list of certnames or domain name globs. Each line represents a node name or group of node names
for which the Autosigner automatically signs certificate requests.

```
hostname.example.com
*.testing.example.com
*.local
```

### AWS and GCP attributes

Using CSR attributes and extensions, a instance deployed to AWS or GCP can embed facts about itself into the certificate request.
Extra details can be found in the Puppet docs: https://puppet.com/docs/puppet/latest/ssl_attributes_extensions.html

The current, supported attributes are:

- pp_cloudplatform
- pp_instance_id
- pp_project
- pp_zone
- pp_region

#### pp_project attribute for AWS Accounts

Puppet doesn't have a `pp_account` extension that we can map to a "AWS Account". For defining the AWS account,
we will reuse the `pp_project` attribute.

#### Creating the CSR with attributes

An example script when deploying a `GCP instance`:
```bash
#!/bin/sh
if [ ! -d /etc/puppetlabs/puppet ]; then
   mkdir /etc/puppetlabs/puppet
fi
cat > /etc/puppetlabs/puppet/csr_attributes.yaml << YAML
extension_requests:
    pp_cloudplatform: gcp
    pp_instance_id: $(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/id)
    pp_project: $(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/project/project-id)
    pp_zone: $(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/zone | cut -d'/' -f4)
YAML
```

An example script when deploying a `AWS instance`.
> Note: The use of `pp_project` for defining the AWS Account number
```bash
#!/bin/sh
if [ ! -d /etc/puppetlabs/puppet ]; then
   mkdir /etc/puppetlabs/puppet
fi
cat > /etc/puppetlabs/puppet/csr_attributes.yaml << YAML
extension_requests:
    pp_cloudplatform: aws
    pp_instance_id: $(curl -s http://169.254.169.254/latest/meta-data/instance-id)
    pp_project: $(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document|grep accountId| awk '{print $3}'|sed  's/"//g'|sed 's/,//g')
    pp_zone: $(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)
YAML
```

test
