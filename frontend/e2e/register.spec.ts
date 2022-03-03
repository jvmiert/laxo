import { test, expect } from "@playwright/test";

test("should navigate to the register page", async ({ page }) => {
  await page.goto("/register");

  await expect(page.locator("//html/body/div/div/main/p")).toContainText(
    "Register",
  );
});
