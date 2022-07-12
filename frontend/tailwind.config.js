module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    fontFamily: {
      sans: ["Inter", "sans-serif"],
      serif: ["serif"],
    },
    extend: {
      animation: {
        pop: "pop 150ms cubic-bezier(0.18, 0.67, 0.6, 1.22)",
      },
      keyframes: {
        pop: {
          "0%": { transform: "translate3d(0px, 0px, 0) scale(1)" },
          "100%": { transform: "translate3d(10px, 10px, 0) scale(1.025)" },
        },
      },
    },
  },
  corePlugins: {
    aspectRatio: false,
  },
  plugins: [require("@tailwindcss/aspect-ratio")],
};
