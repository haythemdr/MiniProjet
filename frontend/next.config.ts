import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "www.tunisianet.com.tn",
      },
      {
        protocol: "https",
        hostname: "tunisianet.com.tn",
      },
      {
        protocol: "https",
        hostname: "wiki.tn",
      },
    ],
  },
};

export default nextConfig;
