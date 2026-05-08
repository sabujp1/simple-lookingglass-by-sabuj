/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: process.env.NEXT_PUBLIC_API_URL + '/:path*',
      },
      {
        source: '/ws',
        destination: process.env.NEXT_PUBLIC_WS_URL + '/ws',
      },
    ];
  },
  images: {
    domains: ['localhost', '127.0.0.1'],
  },
  experimental: {
    serverActions: true,
  },
};

module.exports = nextConfig;