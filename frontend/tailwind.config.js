/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#6C5CE7',
        secondary: '#A29BFE',
        success: '#00B894',
        warning: '#FDCB6E',
        error: '#FF6B6B',
        destructive: '#FF6B6B',
        muted: '#DFE6E9',
        'text-primary': '#2D3436',
        'text-secondary': '#636e72',
      },
      fontFamily: {
        sans: [
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Oxygen',
          'Ubuntu',
          'Cantarell',
          'Fira Sans',
          'Droid Sans',
          'Helvetica Neue',
          'sans-serif',
        ],
      },
    },
  },
  plugins: [],
};
