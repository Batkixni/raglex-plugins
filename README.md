# Raglex Plugins Registry

Official & community plugin marketplace directory for **[Raglex](https://github.com/Batkixni/Raglex)**.

Raglex fetches the live catalog directly from:

```
https://raw.githubusercontent.com/Batkixni/raglex-plugins/main/registry.json
```

Every plugin merged here appears in **Raglex → Preferences → Plugins → Marketplace** for 1-click installation with strict **SHA-256 verification** and explicit permission badges.

---

## Official Plugins

| Plugin | ID | Category | Description |
| :--- | :--- | :---: | :--- |
| **Eagle Import** | `raglex.library-migrator` | `import` | Import Eagle `.library` folders or `.eaglepack` archives into Raglex, restoring nested folders and colors while preserving tags, star ratings, notes, and source URLs. |
| **Smart Duplicate Finder** | `raglex.duplicate-finder` | `utility` | Scan your library for potential duplicate assets by filename and dimensions, automatically tagging matches with `#Duplicate-Review` for quick cleanup. |
| **Markdown Notes & Moodboard Brief** | `raglex.markdown-notes` | `preview` | Preview and edit `.md` notes natively inside Raglex with a borderless writing desk, and export selected reference assets into a structured Markdown brief. |
| **Color Palette Exporter (CSS / JSON)** | `raglex.palette-export` | `export` | Export dominant color palettes from selected assets as CSS custom properties (`:root`) and design tokens ready to paste into frontend projects. |
| **Auto Aspect Ratio & Resolution Tagger** | `raglex.auto-aspect-tagger` | `utility` | Automatically tag images by orientation and resolution (`Portrait`, `Landscape`, `Square`, `Ultrawide`, `4K+`). |

---

## Publishing a Community Plugin (GitHub PR)

Plugins can be written in **any language** (Go, Python, Node.js, Rust, PowerShell, Bash) using Raglex's simple stdin/stdout JSON protocol. Read the full [Raglex Plugin SDK Guide](https://github.com/Batkixni/Raglex/blob/main/docs/PLUGIN_SDK.md) to build and test your plugin locally first.

### Step 1 — Package & Release a `.zip` Archive

Create a GitHub Release on your plugin repository containing a `.zip` archive with `manifest.json` at its root (or inside a single top-level directory):

```text
my-plugin.zip
├── manifest.json
└── run.py          (or compiled binary / script)
```

### Step 2 — Compute the SHA-256 Checksum

- **macOS / Linux**:
  ```bash
  sha256sum my-plugin.zip
  ```
- **Windows (PowerShell)**:
  ```powershell
  (Get-FileHash .\my-plugin.zip -Algorithm SHA256).Hash.ToLower()
  ```

### Step 3 — Submit a Pull Request

1. Fork this repository (`Batkixni/raglex-plugins`).
2. Add a single file at `plugins/<your-plugin-id>/plugin.json` (do **not** edit `registry.json` manually — GitHub Actions rebuilds it automatically on merge):

```json
{
  "id": "yourname.my-plugin",
  "name": "My Awesome Plugin",
  "version": "1.0.0",
  "author": "Your Name",
  "description": "Concise one-sentence English description of what your plugin does.",
  "icon": "palette",
  "accent": "#FF5C33",
  "category": "export",
  "repository": "https://github.com/yourname/raglex-my-plugin",
  "download_url": "https://github.com/yourname/raglex-my-plugin/releases/download/v1.0.0/my-plugin.zip",
  "sha256": "<64-hex-sha256-of-zip>",
  "min_raglex_version": "0.3.3",
  "hooks": ["context_menu", "batch_action"],
  "permissions": ["read_library", "disk"]
}
```

3. Open a Pull Request against `main`.

### Step 4 — Automated CI Validation & Publishing

GitHub Actions (`.github/workflows/validate-and-build.yml`) automatically checks every PR:
- Validates `id` format (`^[a-z0-9][a-z0-9._-]{1,63}$`) and checks for duplicates.
- Verifies required fields (`name`, `version`, `description`, `https://` `download_url`, 64-char `sha256`).
- Downloads the `.zip` archive and verifies its SHA-256 digest matches bit-for-bit.
- Once merged to `main`, CI automatically rebuilds and commits `registry.json`, making your plugin live inside Raglex immediately.

---

## `plugin.json` Field Reference

| Field | Required | Description |
| :--- | :---: | :--- |
| `id` | Yes | Unique lowercase identifier, e.g. `yourname.my-plugin` (regex `^[a-z0-9][a-z0-9._-]{1,63}$`). |
| `name` | Yes | Human-readable plugin title (English recommended for global discovery). |
| `version` | Yes | Semantic version string (e.g. `1.0.0`). Increment this whenever publishing updates so Raglex shows an **Update** button. |
| `author` | Yes | Author name or GitHub handle. |
| `description` | Yes | Clear one-sentence English description shown in the Marketplace list. |
| `download_url` | Yes | Direct `https://` URL to the release `.zip` archive (must be under 64 MB). |
| `sha256` | Yes | 64-character lowercase hex SHA-256 checksum of the `.zip` archive. |
| `repository` | Recommended | Source repository URL for transparency and issue tracking. |
| `category` | Recommended | One of: `import`, `export`, `preview`, `utility`, `ai`. |
| `icon` | Optional | Semantic icon name: `import`, `duplicate`, `markdown`, `palette`, `crop`, `layers`, `folder`, or omit for default sparkle icon. |
| `accent` | Optional | Hex tint color for the icon badge (e.g. `#FF5C33`, `#10B981`, `#38BDF8`). |
| `min_raglex_version` | Optional | Minimum Raglex version required (e.g. `0.4.0`). |
| `hooks` | Optional | Trigger points: `context_menu`, `batch_action`, `library_action`, `importer`, `preview_renderer`. |
| `permissions` | Optional | Declared permissions shown to users: `read_library`, `write_tags`, `disk`, `network`. |

---

## Updating an Existing Plugin

To release a new version of your plugin, open a PR updating `version`, `download_url`, and `sha256` in `plugins/<your-plugin-id>/plugin.json`. Existing users will automatically see an **Update** button in Raglex's Marketplace.

---

## Security & Review Policy

- Every community plugin must use an `https://` download URL with a verified SHA-256 checksum.
- Obfuscated binaries, hidden telemetry, credential harvesting, or undeclared network/file access will be rejected.
- Open-source repositories are strongly preferred so users can audit what runs on their machine.
