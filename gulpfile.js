const autoprefixer = require('autoprefixer');
const browserSync  = require('browser-sync');
const changed      = require('gulp-changed');
const del          = require('del');
const exec         = require('child_process').exec;
const gulp         = require('gulp');
const jasmine      = require('gulp-jasmine-phantom');
const historyApi   = require('connect-history-api-fallback');
const postcss      = require('gulp-postcss');
const sass         = require('gulp-sass');
const sourcemaps   = require('gulp-sourcemaps');
const tslint       = require('gulp-tslint');
const typescript   = require('gulp-typescript');

//=========================================================
//  PATHS
//---------------------------------------------------------
const paths = {
  lib: {
    src: [
      'node_modules/angular2/bundles/angular2.{js,min.js}',
      'node_modules/angular2/bundles/angular2-polyfills.{js,min.js}',
      'node_modules/angular2/bundles/http.{js,min.js}',
      'node_modules/angular2/bundles/router.{js,min.js}',
      'node_modules/rxjs/bundles/Rx.{js,min.js,min.js.map}',
      'node_modules/es6-shim/es6-shim.{js,min.js,min.js.map}',
      'node_modules/angular2/bundles/angular2.dev.{js,min.js,min.js.map}',
      'node_modules/systemjs/dist/system.src.{js,min.js,min.js.map}',
      'node_modules/systemjs/dist/system.{js,js.map}'
    ],
    target: 'webapp/app/lib'
  },

  src: {
    html: 'webapp/**/*.html',
    sass: 'webapp/src/**/*.scss',
    ts: 'webapp/src/**/*.ts'
  },

  target: 'webapp/app',

  typings: {
    entries: 'typings/tsd.d.ts',
    watch: 'typings/**/*.ts'
  }
};

//=========================================================
//  CONFIG
//---------------------------------------------------------
const config = {
  autoprefixer: {
    browsers: ['last 3 versions', 'Firefox ESR']
  },

  browserSync: {
    files: [paths.target + '/**/*'],
    notify: false,
    open: false,
    port: 3000,
    reloadDelay: 500,
    server: {
      baseDir: paths.target,
      middleware: [
        historyApi()
      ]
    }
  },

  sass: {
    errLogToConsole: true,
    outputStyle: 'nested',
    precision: 10,
    sourceComments: false
  },

  ts: {
    configFile: 'tsconfig.json'
  },

  tslint: {
    report: {
      options: {emitError: true},
      type: 'verbose'
    }
  }
};

//=========================================================
//  TASKS
//---------------------------------------------------------
gulp.task('clean.target', () => del(paths.target));

gulp.task('copy.html', () => {
  return gulp.src(paths.src.html)
    .pipe(gulp.dest(paths.target));
});

gulp.task('copy.lib', () => {
  return gulp.src(paths.lib.src)
    .pipe(gulp.dest(paths.lib.target));
});

gulp.task('lint', () => {
  return gulp.src(paths.src.ts)
    .pipe(tslint())
    .pipe(tslint.report(
      config.tslint.report.type,
      config.tslint.report.options
    ));
});

gulp.task('sass', () => {
  return gulp.src(paths.src.sass)
    .pipe(sass(config.sass))
    .pipe(postcss([
      autoprefixer(config.autoprefixer)
    ]))
    .pipe(gulp.dest(paths.target));
});

gulp.task('serve', done => {
  browserSync.create()
    .init(config.browserSync, done);
});

const tsProject = typescript.createProject(config.ts.configFile);

gulp.task('ts', () => {
  return gulp.src([paths.src.ts, paths.typings.entries], {allowEmpty: true})
    .pipe(changed(paths.target, {extension: '.js'}))
    .pipe(sourcemaps.init())
    .pipe(typescript(tsProject))
    .js
    .pipe(sourcemaps.write('.'))
    .pipe(gulp.dest(paths.target));
});

//===========================
//  BUILD
//---------------------------
gulp.task('build', gulp.series(
  'clean.target',
  'copy.html',
  'copy.lib',
  'sass',
  'ts'
));

gulp.task('heroku:production', gulp.series(
    'build'
));

//===========================
//  DEVELOP
//---------------------------
gulp.task('default', gulp.series(
  'build',
  'serve',
  function watch(){
    gulp.watch(paths.src.html, gulp.task('copy.html'));
    gulp.watch(paths.src.sass, gulp.task('sass'));
    gulp.watch([paths.src.ts, paths.typings.watch], gulp.task('ts'));
  }
));

//===========================
//  TEST
//---------------------------

gulp.task('jasmine', function () {
    return gulp.src('webapp/tests/test.js')
        .pipe(jasmine());
});

gulp.task('test', gulp.series('lint', 'build', 'jasmine'));
