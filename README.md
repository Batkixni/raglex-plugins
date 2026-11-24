# Raglex Official & Community Plugin Registry

This repository hosts the plugin catalog (`registry.json`) fetched by **Raglex → Preferences → Plugins → Marketplace**.

## Submitting a Plugin

1. Fork this repository and add `plugins/<your-plugin-id>/plugin.json`.
2. Make sure `download_url` points to an `https://` `.zip` release archive and `sha256` matches the 64-character SHA-256 digest of that archive.
3. Open a Pull Request. GitHub Actions validates every `plugin.json`, verifies the download checksum, and rebuilds `registry.json` automatically upon merge.

See the [Raglex Plugin SDK documentation](https://github.com/Batkixni/Raglex/blob/main/docs/PLUGIN_SDK.md) for manifest and hook specifications.
