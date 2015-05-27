var gulp = require('gulp');
var coffee = require('gulp-coffee');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var imagemin = require('gulp-imagemin');
var sourcemaps = require('gulp-sourcemaps');
var htmlmin = require('gulp-htmlmin');
var minifyCss = require('gulp-minify-css');
var connect = require('gulp-connect');
var revappend = require('gulp-version-append');
var del = require('del');

var paths = {
  public: 'public',
  scripts_3rd: [
    'node_modules/jquery/dist/jquery.js',
    'node_modules/bootstrap/dist/js/bootstrap.js'
  ],
  scripts: 'javascripts/**/*.coffee',
  styles: [
    'node_modules/bootstrap/dist/css/bootstrap.css',
    'stylesheets/**/*.css'
  ],
  templates: 'templates/**/*.html',
  images: [
    'node_modules/bootstrap/dist/fonts/*',
    'images/**/*'
  ]
};

gulp.task('clean', function(cb) {
  del([paths.public], cb);
});

gulp.task('styles', function() {
  return gulp.src(paths.styles)
    .pipe(sourcemaps.init())
    //.pipe(minifyCss({compatibility: 'ie8'}))
    .pipe(minifyCss())
    .pipe(concat('all.min.css'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.public));
});

gulp.task('scripts_3rd', function() {
  return gulp.src(paths.scripts_3rd)
    .pipe(sourcemaps.init())
    .pipe(uglify())
    .pipe(concat('3rd.min.js'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.public));
});

gulp.task('scripts', function() {
  return gulp.src(paths.scripts)
    .pipe(sourcemaps.init())
    .pipe(coffee())
    .pipe(uglify())
    .pipe(concat('main.min.js'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.public));
});

gulp.task('images', function() {
  return gulp.src(paths.images)
    .pipe(imagemin({optimizationLevel: 5}))
    .pipe(gulp.dest(paths.public+'/images'));
});

gulp.task('templates', function() {
  return gulp.src(paths.templates)
    .pipe(revappend(['html', 'js', 'css']))
    .pipe(htmlmin({collapseWhitespace: true}))
    .pipe(gulp.dest(paths.public))
});

gulp.task('watch', function() {
  gulp.watch(paths.scripts, ['scripts']);
  gulp.watch(paths.styles, ['styles']);
  gulp.watch(paths.templates, ['templates']);
  gulp.watch(paths.images, ['images']);
});

gulp.task('server', function() {
  connect.server({
       root: paths.public,
       //livereload: true,
       //fallback: 'index.html',
       port: 8000
   });
});

gulp.task('default', ['watch', 'build']);

gulp.task('build', ['templates', 'scripts', 'scripts_3rd', 'styles', 'images']);
