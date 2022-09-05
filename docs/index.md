# Frontegg + Steampipe

## Documentation

- **[Table definitions & examples →](/plugins/necaris/frontegg/tables)**

## Get started

### Install

Download and install the latest Frontegg plugin:

```bash
steampipe plugin install frontegg
```

### Credentials

| Item | Description |
| - | - |
| Client ID | Zendesk requires an [API token](https://support.zendesk.com/hc/en-us/articles/226022787-Generating-a-new-API-token-), subdomain and email for all requests. |
| Secret | You must be an administrator of your domain to create an API token. |
| Radius | A Frontegg connection is scoped to a single vendor account, with a single set of credentials. |
| Resolution |  1. Credentials specified in environment variables e.g. `FRONTEGG_CLIENT_ID`.<br />2. Credentials in the Steampipe configuration file (`~/.steampipe/config/frontegg.spc`) |

### Configuration

Installing the latest Frontegg plugin will create a config file (`~/.steampipe/config/frontegg.spc`) with a single connection named `frontegg`:

```hcl
connection "frontegg" {
  plugin = "frontegg"
  clientId  = "foo"
  secret    = "bar"
}
```

- `clientId` - The client ID of your Frontegg token
- `secret` - The secret of your Frontegg token

