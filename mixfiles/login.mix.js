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

// =======================
// global
mix.version();
mix.disableNotifications();

mix.babelConfig({
    "presets": [["es2015", {"modules": false}]],
})

var path = './resources/public';
if (!mix.inProduction()) {
    path = './resources/public/dev';
} else {
    //vendors.push('raven-js')
}

var version = "v6"
mix.styles([
    'resources/assets/m/sass/login.css',
], path + '/m/css/' + version + '/login.css');

mix.js('resources/assets/m/js/login.js', path + '/m/js/'+version);

mix.setPublicPath(path + '/m')
mix.setResourceRoot('/m/')