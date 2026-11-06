import { defineConfig, devices } from "@playwright/test";

// Browser tests for the main business, creator and admin flows (e2e/).
// They expect the site and API to be running with e2e/fixtures.sql loaded;
// CI does that on a throwaway database (see .github/workflows/ci.yml).
export default defineConfig({
  testDir: "./e2e",
  // The flows share one database, so they run one at a time.
  workers: 1,
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  timeout: 60_000,
  expect: { timeout: 15_000 },
  reporter: process.env.CI ? [["list"], ["html", { open: "never" }]] : "list",
  use: {
    baseURL: process.env.E2E_BASE_URL ?? "http://localhost:3000",
    viewport: { width: 1440, height: 900 },
    screenshot: "only-on-failure",
    trace: "retain-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        viewport: { width: 1440, height: 900 },
        // Locally, E2E_CHROME=chrome uses the installed Google Chrome instead of downloading Chromium.
        ...(process.env.E2E_CHROME ? { channel: process.env.E2E_CHROME } : {}),
      },
    },
  ],
});
