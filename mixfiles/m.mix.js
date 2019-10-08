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
    resolve: {
        alias: {
            'vue-router$': 'vue-router/dist/vue-router.min.js',
            'vue$': 'vue/dist/vue.min.js'
        }
    }
});

// =======================
// global
let vendors = ['vue', 'axios', 'vuex', 'vue-awesome-swiper'];
if (mix.inProduction()) {
    //vendors.push('raven-js');
}
mix.version();
mix.disableNotifications();
mix.extract(vendors);

mix.babelConfig({
    "presets": [["es2015", {"modules": false}]],
    // "plugins": [
    //     [
    //         "component",
    //         {
    //             "libraryName": "element-ui",
    //             "styleLibraryName": "theme-chalk"
    //         }
    //     ]
    // ]
})

var path = './resources/public';
if (!mix.inProduction()) {
    path = './resources/public/dev';
} else {
    //vendors.push('raven-js')
}

// =======================
// m
var version = "v6"
mix.styles([
    'resources/assets/m/sass/iconfont.css',
    'resources/assets/m/sass/app.scss',
], path + '/m/css/' + version + '/app.css');

mix.js('resources/assets/m/js/app.js', path + '/m/js/'+version);

mix.setPublicPath(path + '/m')
mix.setResourceRoot('/m/')