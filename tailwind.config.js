/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./public/view/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}

