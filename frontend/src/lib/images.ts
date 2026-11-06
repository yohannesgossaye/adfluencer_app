/**
 * Only our own uploads (and files in /public) go through Next.js image
 * resizing; next.config.ts allows just the API host's /uploads path. Anything
 * else (an old link, a local preview of a file being chosen) is shown as it
 * is, so an unknown address can never break a page.
 */
const apiOrigin = (() => {
  try {
    return process.env.NEXT_PUBLIC_API_URL ? new URL(process.env.NEXT_PUBLIC_API_URL).origin : null;
  } catch {
    return null;
  }
})();

export function canOptimize(src?: string | null): boolean {
  if (!src) return false;
  if (src.startsWith("/") && !src.startsWith("//")) return true;
  try {
    const url = new URL(src);
    return apiOrigin !== null && url.origin === apiOrigin && url.pathname.startsWith("/uploads/");
  } catch {
    return false;
  }
}
