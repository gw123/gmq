let mix = require('laravel-mix')
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
let vendors = ['vue', 'axios' ,'vuex'];
if (mix.inProduction()) {
    //vendors.push('raven-js');
}
mix.version();
mix.disableNotifications();
mix.extract(vendors);


// =======================
// pc

// mix.options({
//     extractVueStyles: 'css/vue-styles.css',
//     publicPath: 'public/pc',
//     resourceRoot: '/pc/'
// });

var path = 'public';
if (!mix.inProduction()) {
    path = 'public/dev';
}else{
    //vendors.push('raven-js')
}

mix.sass('resources/assets/pc/sass/app.scss', path + '/pc/css');
mix.styles(['resources/assets/pc/sass/iconfont.css'], path + '/pc/css/iconfont.css');

mix.js('resources/assets/pc/js/app.js', path + '/pc/js').sourceMaps();

mix.setPublicPath(path + '/pc')
mix.setResourceRoot('/pc/')