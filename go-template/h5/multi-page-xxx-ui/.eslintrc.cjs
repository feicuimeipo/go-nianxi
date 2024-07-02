module.exports = {
    root: true,
    "env": {
        "browser": true,
        "es2021": true,
        "node": true,
        "jquery": true,
    },
    "extends": [
        "eslint:recommended",
        "plugin:vue/vue3-essential",
        "plugin:@typescript-eslint/recommended",
        '@vue/prettier',
        '@vue/prettier/@typescript-eslint',
        'plugin:prettier/recommended'
    ],
    "overrides": [
    ],
    "parser":"vue-eslint-parser",
    "parserOptions": {
        "ecmaVersion": "latest",
        //ecmaVersion: 2020,
        "parser": "@typescript-eslint/parser",
        "sourceType": "module"
    },
    "plugins": [
        "vue",
        "@typescript-eslint"
    ],
    "rules": {
        //环境添加
        'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
        'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
        //"Parsing error: x-invalid-end-tag"
        //"vue/no-parsing-error": [2, { "x-invalid-end-tag": false }],
        //关闭组件命名规则-vue/multi-word-component-names
        "vue/multi-word-component-names":"off",
        "@typescript-eslint/ban-types": [
            "error",{
                "extendDefaults": true,
                "types": {
                    "{}": false
                }
            }
        ]
    }
}
