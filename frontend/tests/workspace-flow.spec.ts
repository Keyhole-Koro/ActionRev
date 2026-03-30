import { expect, test } from '@playwright/test'

test('can upload a document into a workspace and see it become ready', async ({ page }) => {
  await page.goto('/dashboard')

  const workspaceLink = page.getByRole('link', { name: /Growth Strategy Review/i })
  await expect(workspaceLink).toBeVisible({ timeout: 15_000 })
  await workspaceLink.click()

  const uploadToggle = page.getByRole('button', { name: 'Upload file' })
  await expect(uploadToggle).toBeVisible({ timeout: 15_000 })

  await uploadToggle.click()
  await expect(page.getByRole('button', { name: 'Upload and process' })).toBeVisible()

  await page.locator('input[type="file"]').setInputFiles({
    name: 'playwright-upload.txt',
    mimeType: 'text/plain',
    buffer: Buffer.from('Synthify Playwright upload fixture'),
  })

  await page.getByRole('button', { name: 'Upload and process' }).click()

  await expect(page.getByText(/nodes · .* edges/)).toBeVisible()

  await page.getByRole('button', { name: 'Upload file' }).click()
  const uploadedDocument = page.getByText('playwright-upload.txt').last()
  await expect(uploadedDocument).toBeVisible()
  await expect(uploadedDocument.locator('..')).toContainText('Ready', { timeout: 10_000 })
})
