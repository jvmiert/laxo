const { PHASE_DEVELOPMENT_SERVER } = require("next/constants");

/** @type {import('next').NextConfig} */
const withBundleAnalyzer = require("@next/bundle-analyzer")({
  enabled: process.env.ANALYZE === "true",
});

module.exports = (phase) =>
  withBundleAnalyzer({
    reactStrictMode: true,
    poweredByHeader: false,
    i18n: {
      locales: ["default", "en", "vi", "sw"],
      defaultLocale: "default",
      localeDetection: false,
    },
    trailingSlash: false,
    async rewrites() {
      return [
        {
          source: "/api/:path*",
          destination: "http://127.0.0.1:8080/api/:path*",
        },
      ];
    },
    env: {
      // Enable signup only for development, disable for production
      ENABLE_SIGN_UP: phase === PHASE_DEVELOPMENT_SERVER ? true : false,
    },
    webpack(config, { dev, ...other }) {
      if (!dev) {
        // https://formatjs.io/docs/guides/advanced-usage#react-intl-without-parser-40-smaller
        config.resolve.alias["@formatjs/icu-messageformat-parser"] =
          "@formatjs/icu-messageformat-parser/no-parser";
      }
      return config;
    },
  });
