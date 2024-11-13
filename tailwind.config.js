/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/views/**/*.templ'
  ],
  theme: {
    container: {
      center: true,
      padding: {
        DEFAULT: "1rem",
        mobile: "2rem",
        tablet: "4rem",
        desktop: "5rem",
      },
    },
    extend: {
    },
  },
  plugins: [
    require('daisyui')
  ],
}
