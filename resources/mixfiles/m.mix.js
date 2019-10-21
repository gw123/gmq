let mix = require('laravel-mix');

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
//let vendors = ['vue', 'axios', 'vuex'];
let vendors = ['vue', 'axios'];
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
});


var path = './public';
if (!mix.inProduction()) {
    path = './public/dev';
} else {
    //vendors.push('raven-js')
}

// =======================
// m
var srcPath = './assets/m';
var version = "v7";
mix.styles([
    srcPath + '/sass/iconfont.css',
    srcPath + '/sass/app.scss',
], path + '/m/css/' + version + '/app.css');

mix.js(srcPath + '/js/app.js', path + '/m/js/' + version);
mix.js(srcPath + '/js/news.js', path + '/m/js/' + version);
mix.js(srcPath + '/js/group.js', path + '/m/js/' + version);
mix.js(srcPath + '/js/home.js', path + '/m/js/' + version);
mix.js(srcPath + '/js/tagNews.js', path + '/m/js/' + version);

//登录页面
mix.styles([
    srcPath + '/sass/login.css',
], path + '/m/css/' + version + '/login.css');

mix.js(srcPath + '/js/login.js', path + '/m/js/'+version);

mix.setPublicPath(path + '/m');
mix.setResourceRoot('/m/');