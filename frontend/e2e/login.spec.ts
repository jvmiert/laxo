import { test, expect } from "@playwright/test";

test("should navigate to the login page", async ({ page }) => {
  await page.goto("/login");

  await expect(page.locator("//html/body/div/div/main/p")).toContainText(
    "Login",
  );
});
