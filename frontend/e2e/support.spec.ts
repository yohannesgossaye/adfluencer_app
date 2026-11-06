import { expect, test } from "@playwright/test";

import { pageAs } from "./helpers";

const SUBJECT = "E2E: payout question";

test("creator asks support, the admin replies and the creator answers", async ({ browser }) => {
  const creator = await pageAs(browser, "creator");
  await creator.goto("/influencer/support");
  await creator.locator("#support-subject").fill(SUBJECT);
  await creator.locator("#support-message").fill("When will I receive my payout for the tasting video?");
  await creator.getByRole("button", { name: "Send request" }).click();
  await expect(creator.getByRole("button", { name: new RegExp(SUBJECT) })).toBeVisible();

  const admin = await pageAs(browser, "admin");
  await admin.goto("/admin/support");
  await admin.locator("article").filter({ hasText: SUBJECT }).getByRole("button", { name: "View details" }).click();
  await admin.getByLabel("Reply to the user").fill("Payouts are sent within 2 working days after approval.");
  await admin.getByRole("button", { name: "Send reply" }).click();
  await expect(admin.getByText("The user has been emailed")).toBeVisible();

  await creator.reload();
  const ticket = creator.getByRole("button", { name: new RegExp(SUBJECT) });
  await expect(ticket.getByText("Support replied")).toBeVisible();
  await ticket.click();
  await expect(creator.getByText("Payouts are sent within 2 working days after approval.")).toBeVisible();
  await creator.getByLabel("Add a message").fill("Thank you!");
  await creator.getByRole("button", { name: "Send", exact: true }).click();
  await expect(creator.getByText("Reply sent")).toBeVisible();
});
