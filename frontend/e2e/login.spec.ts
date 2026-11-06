import { test } from "@playwright/test";

import { loginAs, type Role } from "./helpers";

for (const role of ["business", "creator", "admin"] as Role[]) {
  test(`${role} logs in and reaches their dashboard`, async ({ page }) => {
    await loginAs(page, role);
  });
}
