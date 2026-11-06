import { expect, test, type Page } from "@playwright/test";

import { loginAs } from "./helpers";

const CREATOR = "e2e00000-0000-4000-8000-000000000002";
const CAMPAIGN = "e2e00000-0000-4000-8000-0000000000c2";

async function openChatFrom(page: Page, profileUrl: string) {
  await page.goto(profileUrl);
  await page.getByRole("button", { name: "Direct Message" }).click();
  await expect(page).toHaveURL(/conversation_id=/);
  return new URL(page.url()).searchParams.get("conversation_id");
}

// One business and one creator have one conversation, whichever page the
// chat is opened from (with a campaign or without, e.g. from a contract).
test("opening a chat from different pages reuses the same conversation", async ({ page }) => {
  await loginAs(page, "business");
  const fromCampaign = await openChatFrom(page, `/business/influencers/${CREATOR}?campaign_id=${CAMPAIGN}`);
  const withoutCampaign = await openChatFrom(page, `/business/influencers/${CREATOR}`);
  expect(withoutCampaign).toBe(fromCampaign);
  await expect(page.getByRole("button", { name: /E2E Creator/ })).toHaveCount(1);
});
