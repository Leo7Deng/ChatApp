import type { NextConfig } from "next";

const isProd = process.env.NODE_ENV === "production";

const nextConfig = {
  reactStrictMode: true,
  images: {
    unoptimized: true,
  },
  assetPrefix: isProd ? '/ChatApp' : '',
  basePath: isProd ? '/ChatApp' : '',
  output: 'export'
};

export default nextConfig;
