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

var src = "assets/m"
var path = './public';
if (!mix.inProduction()) {
    path = './public/dev';
} else {
    //vendors.push('raven-js')
}

var version = "v6"
mix.styles([
    src + '/sass/login.css',
], path + '/m/css/' + version + '/login.css');

mix.js(src + '/js/login.js', path + '/m/js/'+version);

mix.setPublicPath(path + '/m')
mix.setResourceRoot('/m/')