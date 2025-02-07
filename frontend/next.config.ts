import type { NextConfig } from "next";

const isProd = process.env.NODE_ENV === 'production';

const nextConfig = {
  reactStrictMode: true,
  images: {
    unoptimized: true,
  },
  assetPrefix: isProd ? '/chatapp' : '',
  basePath: isProd ? '/chatapp' : '',
  output: 'export'
};

export default nextConfig;
