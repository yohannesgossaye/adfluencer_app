import { expect, test } from "@playwright/test";

import { loginAs } from "./helpers";

// Bank-style sessions: a warning 60 s before the idle limit, "Stay signed in"
// keeps you in, and doing nothing signs you out to the login page.
// The page clock is fast-forwarded instead of waiting 15 real minutes.
test("inactive business is warned, can stay, then is signed out", async ({ page }) => {
  await page.clock.install();
  await loginAs(page, "business");
  const warning = page.getByRole("dialog").filter({ hasText: "Are you still there?" });

  await page.clock.fastForward("13:10");
  await expect(warning).toBeVisible();
  await warning.getByRole("button", { name: "Stay signed in" }).click();
  await expect(warning).toBeHidden();
  await expect(page).toHaveURL(/\/business\/dashboard/);

  await page.clock.fastForward("13:10");
  await expect(warning).toBeVisible();
  await page.clock.fastForward("01:00");
  await expect(page).toHaveURL(/\/auth\/login\?.*reason=idle/);
  await expect(page.getByText("signed out after a period of inactivity")).toBeVisible();
});

test("admin limit is shorter: warning after about 8 minutes", async ({ page }) => {
  await page.clock.install();
  await loginAs(page, "admin");
  const warning = page.getByRole("dialog").filter({ hasText: "Are you still there?" });

  await page.clock.fastForward("07:30");
  await expect(warning).toBeHidden();
  await page.clock.fastForward("00:45");
  await expect(warning).toBeVisible();
});

test("the login page has no Remember me", async ({ page }) => {
  await page.goto("/auth/login");
  await expect(page.locator("#page-login-email")).toBeVisible();
  await expect(page.getByText("Remember me")).toHaveCount(0);
});
