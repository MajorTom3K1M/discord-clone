/** @type {import('next').NextConfig} */
const nextConfig = {
    webpack: (config) => {
      config.externals.push({
        "utf-8-validate": "commonjs utf-8-validate",
        bufferutil: "commonjs bufferutil",
        canvas: "canvas",
      });
  
      return config;
    },
    images: {
      domains: [
        "uploadthing.com",
        "utfs.io"
      ]
    },
    // reactStrictMode: false,
  }

module.exports = nextConfig
