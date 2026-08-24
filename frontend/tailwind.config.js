/** @type {import('tailwindcss').Config} */
export default {
    // HTMLとJSファイルをスキャン対象にする
    content: [
        "../static/js/**/*.js",
        "../templates/**/*.html",
        "../internal/**/*.go",
    ],
    theme: {
        extend: {
            theme: {
                extend: {
                    colors: {
                        ink: "#14213D",
                        navy: {
                            900: "#14213D",
                            800: "#243554",
                            700: "#334A70",
                            600: "#476084",
                            500: "#62789A",
                        },
                        accent: {
                            DEFAULT: "#7A1F32",
                            bright: "#9B4053",
                            deep: "#651929",
                        },
                        sev: { hi: "#C2415A", md: "#B7791F", lo: "#248A69" },
                    },
                    fontFamily: {
                        display: [
                            '"Space Grotesk"',
                            '"Noto Sans JP"',
                            "sans-serif",
                        ],
                        sans: ["Inter", '"Noto Sans JP"', "sans-serif"],
                    },
                },
            },
        },
    },
    plugins: [require("daisyui")],
    // daisyUIの設定（任意）
    daisyui: {
        themes: ["light", "dark", "cupcake", "coffee", "emerald", "abyss"],
    },
};
