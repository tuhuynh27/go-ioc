import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'Go IoC',
  tagline: 'Bring "@Autowired" to Go! Compile-time dependency injection with Spring-like syntax.',
  favicon: 'img/favicon.ico',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: 'https://go-ioc.keva.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'tuhuynh27', // Usually your GitHub org/user name.
  projectName: 'go-ioc', // Usually your repo name.

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/tuhuynh27/go-ioc/tree/main/website/',
        },
        blog: false, // Disable blog for now
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/logo.png',
    navbar: {
      title: 'Go IoC',
      logo: {
        alt: 'Go IoC Logo',
        src: 'img/logo.png',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          to: '/docs/vscode-extension',
          label: 'VS Code Extension',
          position: 'left',
        },
        {
          href: 'https://marketplace.visualstudio.com/items?itemName=keva-dev.go-ioc',
          label: 'VS Code Marketplace',
          position: 'right',
        },
        {
          href: 'https://github.com/tuhuynh27/go-ioc',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Documentation',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/intro',
            },
            {
              label: 'Testing Guide',
              to: '/docs/testing',
            },
            {
              label: 'VS Code Extension',
              to: '/docs/vscode-extension',
            },
          ],
        },
        {
          title: 'Tools',
          items: [
            {
              label: 'VS Code Extension',
              href: 'https://marketplace.visualstudio.com/items?itemName=keva-dev.go-ioc',
            },
            {
              label: 'GitHub Repository',
              href: 'https://github.com/tuhuynh27/go-ioc',
            },
            {
              label: 'Example Project',
              href: 'https://github.com/tuhuynh27/go-ioc-gin-demo',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'Issues',
              href: 'https://github.com/tuhuynh27/go-ioc/issues',
            },
            {
              label: 'GitHub Discussions',
              href: 'https://github.com/tuhuynh27/go-ioc/discussions',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Go IoC. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'bash', 'json'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
