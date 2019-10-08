let mix = require('laravel-mix');
/*
 |--------------------------------------------------------------------------
 | Mix Asset Management
 |--------------------------------------------------------------------------
 |
 | Mix provides a clean, fluent API for defining some Webpack build steps
 | for your Laravel application. By default, we are compiling the Sass
 | file for the application as well as bundling up all the JS files.
 |
 */
mix.webpackConfig({
    resolve:{
        alias: {
            'vue-router$': 'vue-router/dist/vue-router.common.js',
            'vue$': 'vue/dist/vue.min.js'
        }
    }
});

// =======================
// global
let vendors = ['vue', 'axios' ,'vuex','vue-router'];
if (mix.inProduction()) {
    //vendors.push('raven-js');
}
mix.version();
mix.disableNotifications();
mix.extract(vendors);
