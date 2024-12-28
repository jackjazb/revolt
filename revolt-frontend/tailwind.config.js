/** @type {import('tailwindcss').Config} */
export default {
    content: ["./src/**/*.svelte"],
    theme: {
        extend: {},
        fontSize: {
            sm: '0.8rem',
            base: '1.1rem',
            lg: '1.2rem',
            xl: '1.4rem',
            '2xl': '1.563rem',
            '3xl': '1.953rem',
            '4xl': '2.441rem',
            '5xl': '3.052rem',
        },
        fontFamily: {
            'sans': ['"Grenz Gotisch"'],
            'mono': ['monospace']
        }
    },
    plugins: [],
}

