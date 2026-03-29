import type { Config } from 'tailwindcss';

export default {
  darkMode: ['class'],
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    container: {
      center: true,
      padding: '1.5rem',
      screens: {
        '2xl': '1440px',
      },
    },
    extend: {
      fontFamily: {
        sans: [
          '"Avenir Next"',
          '"PingFang SC"',
          '"Hiragino Sans GB"',
          '"Microsoft YaHei"',
          'sans-serif',
        ],
        display: ['"Avenir Next Condensed"', '"Avenir Next"', '"PingFang SC"', 'sans-serif'],
      },
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        card: 'hsl(var(--card))',
        'card-foreground': 'hsl(var(--card-foreground))',
        popover: 'hsl(var(--popover))',
        'popover-foreground': 'hsl(var(--popover-foreground))',
        primary: 'hsl(var(--primary))',
        'primary-foreground': 'hsl(var(--primary-foreground))',
        secondary: 'hsl(var(--secondary))',
        'secondary-foreground': 'hsl(var(--secondary-foreground))',
        muted: 'hsl(var(--muted))',
        'muted-foreground': 'hsl(var(--muted-foreground))',
        accent: 'hsl(var(--accent))',
        'accent-foreground': 'hsl(var(--accent-foreground))',
        destructive: 'hsl(var(--destructive))',
        'destructive-foreground': 'hsl(var(--destructive-foreground))',
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        sidebar: 'hsl(var(--sidebar))',
        'sidebar-foreground': 'hsl(var(--sidebar-foreground))',
        'sidebar-accent': 'hsl(var(--sidebar-accent))',
        'sidebar-accent-foreground': 'hsl(var(--sidebar-accent-foreground))',
        'sidebar-border': 'hsl(var(--sidebar-border))',
      },
      borderRadius: {
        xl: 'calc(var(--radius) + 6px)',
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 6px)',
      },
      boxShadow: {
        soft: '0 18px 48px rgba(15, 23, 42, 0.08)',
        glow: '0 24px 80px rgba(14, 165, 233, 0.18)',
      },
      backgroundImage: {
        'hero-grid':
          'radial-gradient(circle at top left, rgba(20, 184, 166, 0.24), transparent 28%), radial-gradient(circle at top right, rgba(56, 189, 248, 0.18), transparent 28%), linear-gradient(135deg, rgba(15, 23, 42, 0.96) 0%, rgba(17, 24, 39, 0.94) 48%, rgba(15, 118, 110, 0.88) 100%)',
      },
    },
  },
  plugins: [],
} satisfies Config;
