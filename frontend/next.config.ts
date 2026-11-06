import type { NextConfig } from "next";

// Uploaded photos are served by the API under /uploads (see backend
// STORAGE_PUBLIC_URL). Only that address is resized by Next.js, so the image
// optimizer can't be used to fetch arbitrary websites. src/lib/images.ts
// decides per image.
function uploadsPattern() {
  try {
    const api = new URL(process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080");
    return {
      protocol: api.protocol.replace(":", "") as "http" | "https",
      hostname: api.hostname,
      port: api.port,
      pathname: "/uploads/**",
    };
  } catch {
    return null;
  }
}

const uploads = uploadsPattern();

const nextConfig: NextConfig = {
  devIndicators: false,
  turbopack: {
    root: process.cwd(),
  },
  images: {
    remotePatterns: uploads ? [uploads] : [],
    formats: ["image/avif", "image/webp"],
  },
};

export default nextConfig;
