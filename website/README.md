# Go IoC Website

This website is built using [Docusaurus](https://docusaurus.io/), a modern static website generator.

## Installation

```bash
npm install
```

## Local Development

```bash
npm start
```

This command starts a local development server and opens up a browser window. Most changes are reflected live without having to restart the server.

## Build

```bash
npm run build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Deployment

The website is automatically deployed to https://go-ioc.keva.dev when changes are pushed to the main branch.

### Manual Deployment

```bash
npm run deploy
```

If you are using GitHub pages for hosting, this command is a convenient way to build the website and push to the `gh-pages` branch.

## Configuration

The main configuration is in `docusaurus.config.ts`. Key settings:

- **Title**: Go IoC
- **URL**: https://go-ioc.keva.dev
- **Repository**: tuhuynh27/go-ioc

## Content Structure

- `docs/` - Documentation pages
- `src/pages/` - Landing page and custom pages
- `src/components/` - React components
- `static/` - Static assets (images, CNAME, etc.)

## Adding Documentation

1. Create a new `.md` file in the `docs/` directory
2. Add frontmatter with `sidebar_position` for ordering
3. The sidebar will automatically update

Example:
```markdown
---
sidebar_position: 3
---

# Your Page Title

Your content here...
```
