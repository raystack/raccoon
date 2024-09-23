module.exports = {
  docsSidebar: [
    'introduction',
    'quickstart',
    'roadmap',
    {
      type: "category",
      label: "Guides",
      items: [
        "guides/overview",
        "guides/publishing",
        "guides/deployment",
        "guides/monitoring",
        "guides/troubleshooting",
      ],
    },
    {
      type: "category",
      label: "Clients",
      items: [
        "clients/overview",
        "clients/golang",
        "clients/python",
        "clients/javascript",
        "clients/java",
      ],
    },
    {
      type: "category",
      label: "Concepts",
      items: [
        "concepts/architecture",
        "concepts/structure",
      ],
    },
    {
      type: "category",
      label: "Reference",
      items: [
        "reference/configurations",
        "reference/metrics",
      ],
    },
    {
      type: "category",
      label: "Contribute",
      items: [
        "contribute/contribution",
        "contribute/development",
        "contribute/release",
      ],
    },
  ],
};