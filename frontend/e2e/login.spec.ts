import { expect, test } from "@playwright/test";

import { loginAs, type Role } from "./helpers";

for (const role of ["business", "creator", "admin"] as Role[]) {
  test(`${role} logs in and reaches their dashboard`, async ({ page }) => {
    await loginAs(page, role);
  });
}

// Uses an email nobody has, so the shared test accounts are never locked.
test("repeated wrong passwords pause logins and say for how long", async ({ page }) => {
  await page.goto("/auth/login");
  await page.locator("#page-login-email").fill("not.a.member@example.com");
  const submit = page.locator("form:has(#page-login-email) button[type=submit]");
  for (let i = 1; i <= 4; i++) {
    await page.locator("#page-login-password").fill(`wrong-${i}`);
    await submit.click();
    // Same answer as for a real account, so nobody can find out who has one.
    await expect(page.getByText("That email and password don't match")).toBeVisible();
  }
  await page.locator("#page-login-password").fill("wrong-5");
  await submit.click();
  await expect(page.getByText("Too many failed attempts. Try again in 5 minutes, or reset your password.")).toBeVisible();
});
