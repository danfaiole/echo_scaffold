/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './internal/views/**/*.templ'
  ],
  theme: {
  },
  plugins: [
    require('daisyui')
  ],
}
