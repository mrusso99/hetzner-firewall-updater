# Hetzner Firewall Updater

A simple tool to automatically update Hetzner firewall rules when your IP changes. Perfect for users without a static IP who want to restrict access to specific applications without using a VPN.

⚠️ **Note:** This project is for personal use only. Feature requests may not be considered.

## How It Works

1. The application pings a specified host (with your firewall configured to accept pings only from your IP).
2. If the host does not respond, it detects that your IP has changed.
3. The tool then updates the Hetzner firewall rules via API to allow your new IP.

## Usage

You can update multiple firewalls at once by passing their names as parameters:

```bash
firewall-updater firewall-1 firewall-2 firewall-3
```

## Configuration

Create a `.env` file with the following variables:

```
HETZNER_TOKEN=your-hetzner-api-token
HOST=host-to-ping
```

- `HETZNER_TOKEN`: Your Hetzner API token.
- `HOST`: The host you want the application to ping.
