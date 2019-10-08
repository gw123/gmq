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
let vendors = ['vue'];

mix.version();
mix.disableNotifications();
mix.extract(vendors);


// =======================
// game

// mix.options({
//     extractVueStyles: 'css/vue-styles.css',
//     publicPath: 'public/game',
//     resourceRoot: '/game/'
// });

mix.sass('resources/assets/game/sass/app.scss', 'public/game/css');

mix.js('resources/assets/game/js/app.js', 'public/game/js');

mix.setPublicPath('public/game')
mix.setResourceRoot('/game/')