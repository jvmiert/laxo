import { test, expect } from "@playwright/test";

test("should navigate to the landing page", async ({ page }) => {
  await page.goto("/");

  // the page should contain hello world
  await expect(page.locator("//html/body/div/div/main/p[1]")).toHaveText(
    "Hello World",
  );
});
