import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './tests',
  timeout: 30_000,
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'list',
  use: {
    baseURL: 'http://127.0.0.1:4273',
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command:
        'PORT=28080 CORS_ALLOWED_ORIGINS=http://127.0.0.1:4273,http://localhost:4273 PUBLIC_BASE_URL=http://127.0.0.1:28080 go run ./cmd/server',
      url: 'http://127.0.0.1:28080/healthz',
      reuseExistingServer: false,
      cwd: '../backend',
    },
    {
      command:
        'NEXT_PUBLIC_API_BASE_URL=http://127.0.0.1:28080 npx next build && NEXT_PUBLIC_API_BASE_URL=http://127.0.0.1:28080 npx next start --hostname 127.0.0.1 --port 4273',
      url: 'http://127.0.0.1:4273/dashboard',
      reuseExistingServer: false,
      cwd: '.',
    },
  ],
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})
