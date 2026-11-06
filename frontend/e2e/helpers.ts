import { expect, type Browser, type Page } from "@playwright/test";

/** Test accounts from e2e/fixtures.sql. */
export const PASSWORD = "E2e-pass-2026";
export const EMAIL = {
  business: "e2e.business@example.com",
  creator: "e2e.creator@example.com",
  admin: "e2e.admin@example.com",
} as const;
export type Role = keyof typeof EMAIL;

const HOME: Record<Role, RegExp> = {
  business: /\/business\/dashboard/,
  creator: /\/influencer\/dashboard/,
  admin: /\/admin\/dashboard/,
};

/** Logs in through the real login page and waits for the dashboard. */
export async function loginAs(page: Page, role: Role) {
  if (role === "admin") {
    await page.goto("/admin/login");
    await page.locator("#admin-email").fill(EMAIL.admin);
    await page.locator("#admin-password").fill(PASSWORD);
    await page.getByRole("button", { name: "Log in", exact: true }).click();
  } else {
    await page.goto("/auth/login");
    await page.locator("#page-login-email").fill(EMAIL[role]);
    await page.locator("#page-login-password").fill(PASSWORD);
    await page.locator("form:has(#page-login-email) button[type=submit]").click();
  }
  await expect(page).toHaveURL(HOME[role]);
}

/** A fresh, logged-in page for one role, so several people can act in one test. */
export async function pageAs(browser: Browser, role: Role) {
  const page = await (await browser.newContext()).newPage();
  await loginAs(page, role);
  return page;
}
