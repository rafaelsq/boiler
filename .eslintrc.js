module.exports = {
  extends: "eslint:recommended",
  env: {
    browser: true
  },
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
      experimentalObjectRestSpread: true
    },
    sourceType: "module"
  },
  plugins: ["react"],
  rules: {
    "react/jsx-uses-vars": ["error"],

    "no-unused-vars": ["error", { varsIgnorePattern: "^h$" }],
    "max-len": ["error", 120],
    indent: ["error", 4],
    "linebreak-style": ["error", "unix"],
    "array-bracket-spacing": ["error", "never"],
    "object-curly-spacing": ["error", "never"],
    "comma-dangle": ["error", "always-multiline"],
    "no-extra-parens": ["error", "all"],
    quotes: ["error", "single"],
    semi: ["error", "never"],
    "quote-props": ["error", "consistent-as-needed"],
    "no-trailing-spaces": ["error"]
  }
};
