const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

// With JSDoc @type annotations, IDEs can provide config autocompletion
/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  title: 'Raccoon',
  tagline: 'Raccoon is a high-throughput, low-latency service to collect events in real-time from your web, mobile apps, and services using multiple network protocols.',
  url: 'https://goto.github.io/',
  baseUrl: '/raccoon/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'favicon.ico',
  organizationName: 'goto',
  projectName: 'raccoon',

  presets: [
    [
      '@docusaurus/preset-classic',
      ({
        gtag: {
          trackingID: 'G-XXX',
        },
        docs: {
          routeBasePath: '/',
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/goto/raccoon/edit/master/docs/',
          sidebarCollapsed: false,
          breadcrumbs: false,
        },
        blog: false,
        theme: {
          customCss: [
            require.resolve('./src/css/theme.css'),
            require.resolve('./src/css/icons.css'),
          ],
        },
      })
    ],
  ],

  themeConfig:
    ({
      colorMode: {
        defaultMode: 'light',
        respectPrefersColorScheme: true,
      },
      navbar: {
        title: 'Raccoon',
        logo: { src: '/logo.svg', },
        hideOnScroll: true,
        items: [
          {
            type: 'doc',
            docId: 'introduction',
            position: 'right',
            label: 'Documentation',
          },
          { to: '/support', label: 'Support', position: 'right' },
          {
            href: 'https://bit.ly/2RzPbtn',
            position: 'right',
            className: 'header-slack-link',
          },
          {
            href: 'https://github.com/goto/raccoon',
            className: 'navbar-item-github',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'light',
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      announcementBar: {
        id: 'star-repo',
        content: '⭐️ If you like Raccoon, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/goto/raccoon">GitHub</a>! ⭐',
        backgroundColor: '#222',
        textColor: '#eee',
        isCloseable: true,
      },
    }),
});
