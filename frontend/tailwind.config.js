/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./generator.html",
    "./builder.html",
    "./assets/js/**/*.js",
    "../../api/*.go",
    "../../internal/builder/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require("@tailwindcss/typography"),
    require("daisyui"),
  ],
  daisyui: {
    themes: ["dark", "corporate"],
    darkTheme: "dark",
  },
}
