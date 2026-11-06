import { expect, test } from "@playwright/test";

import { loginAs } from "./helpers";

test("creator accepts a campaign invite", async ({ page }) => {
  await loginAs(page, "creator");
  await page.goto("/influencer/invites");

  const invite = page.locator("article").filter({ hasText: "E2E food tasting" });
  await expect(invite.getByText("Needs your answer")).toBeVisible();
  await invite.getByRole("button", { name: "Accept" }).click();

  await expect(page.getByText("Invite accepted")).toBeVisible();
  await expect(invite.getByText("Accepted", { exact: true })).toBeVisible();
});
