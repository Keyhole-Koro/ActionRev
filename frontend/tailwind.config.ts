import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './features/**/*.{js,ts,jsx,tsx,mdx}',
    './lib/**/*.{js,ts,jsx,tsx,mdx}',
    './src/generated/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-inter)', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        ink: '#0f172a',
        mist: '#e2e8f0',
        panel: '#f8fafc',
        accent: '#0f766e',
        accentSoft: '#ccfbf1',
      },
      boxShadow: {
        panel: '0 1px 3px rgba(15,23,42,0.06), 0 4px 16px rgba(15,23,42,0.06)',
        card: '0 1px 2px rgba(15,23,42,0.04), 0 8px 32px rgba(15,23,42,0.08)',
      },
    },
  },
  plugins: [],
}

export default config
