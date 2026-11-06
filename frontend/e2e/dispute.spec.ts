import { expect, test } from "@playwright/test";

import { pageAs } from "./helpers";

const CONTRACT = "/business/contracts?contract_id=e2e00000-0000-4000-8000-0000000000d1";

test("business reports a problem on a contract and the admin resolves it", async ({ browser }) => {
  // The business reports a problem on the paid, in-progress contract.
  const business = await pageAs(browser, "business");
  await business.goto(CONTRACT);
  await business.getByRole("button", { name: "Report a problem" }).click();
  const report = business.getByRole("dialog");
  await report.getByLabel("Work not delivered").check();
  await report.getByLabel("Tell us what happened").fill("The tasting video was due on Monday and has not been posted.");
  await report.getByRole("button", { name: "Report problem" }).click();
  await expect(business.getByText("Problem reported: Work not delivered")).toBeVisible();

  // The admin decides: the contract continues.
  const admin = await pageAs(browser, "admin");
  await admin.goto("/admin/disputes");
  await expect(admin.getByText("Reported by:")).toBeVisible();
  await admin.getByRole("button", { name: /Resolve dispute/ }).click();
  const resolve = admin.getByRole("dialog");
  await resolve.getByLabel("Continue the contract").check();
  await resolve.getByLabel("Note to both sides").fill("The creator will post the video today.");
  await resolve.getByRole("button", { name: "Continue the contract" }).click();
  await expect(admin.getByText("Both sides have been emailed")).toBeVisible();

  // The business sees the decision on the contract.
  await business.goto(CONTRACT);
  await expect(business.getByText("Problem resolved: Continue the contract")).toBeVisible();
});
