import { expect, test } from "@playwright/test";

import { loginAs } from "./helpers";

test("business searches creators, ranks them for a campaign and invites one", async ({ page }) => {
  await loginAs(page, "business");
  await page.goto("/business/recommendations");

  // Works without a campaign: all approved creators, then a search.
  await expect(page.getByText(/\d+ creators/)).toBeVisible();
  await page.getByLabel("Search creators").fill("football");
  const card = page.getByRole("listitem").filter({ hasText: "Bereket Football Talk" });
  await expect(card).toBeVisible();
  await expect(page.getByText("2 creators")).toBeVisible();

  // Ranked for the live sports campaign, the creator can be invited.
  await page.getByLabel("Rank for campaign").selectOption({ label: "E2E campaign" });
  await expect(page.getByText("ranked for")).toBeVisible();
  await expect(card.getByText("Campaign match")).toBeVisible();
  await card.getByRole("button", { name: "Invite" }).click();
  await expect(card.getByRole("button", { name: "Invited" })).toBeVisible();
});
