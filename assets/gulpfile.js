var gulp = require('gulp');
var coffee = require('gulp-coffee');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var imagemin = require('gulp-imagemin');
var sourcemaps = require('gulp-sourcemaps');
//var minifyHTML = require('gulp-minify-html');
var htmlmin = require('gulp-htmlmin');
// var webserver = require('gulp-webserver');
var connect = require('gulp-connect');
var del = require('del');

var paths = {
  scripts: 'javascripts/**/*.coffee',
  templates: 'templates/**/*.html',
  images: 'images/**/*'
};

gulp.task('clean', function(cb) {
  del(['build'], cb);
});

gulp.task('scripts', ['clean'], function() {
  return gulp.src(paths.scripts)
    .pipe(sourcemaps.init())
      .pipe(coffee())
      .pipe(uglify())
      .pipe(concat('all.min.js'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('build/javascripts'));
});

gulp.task('images', ['clean'], function() {
  return gulp.src(paths.images)
    .pipe(imagemin({optimizationLevel: 5}))
    .pipe(gulp.dest('build/images'));
});

gulp.task('templates', function() {
  return gulp.src(paths.templates)
    .pipe(htmlmin({collapseWhitespace: true}))
    .pipe(gulp.dest('build'))

  // var opts = {
  //   conditionals: true,
  //   spare:true
  // };
  //
  // return gulp.src(paths.templates)
  //   .pipe(minifyHTML(opts))
  //   .pipe(gulp.dest('build'));
});

gulp.task('watch', function() {
  gulp.watch(paths.scripts, ['scripts']);
  gulp.watch(paths.templates, ['templates']);
  gulp.watch(paths.images, ['images']);
});

gulp.task('server', function() {
  connect.server({
       root: 'build',
       //livereload: true,
       //fallback: 'index.html',
       port: 8000
   });

  // gulp.src('build')
  //   .pipe(webserver({
  //     livereload: true,
  //     directoryListing: true,
  //     open: true
  //   }));
});

gulp.task('default', ['clean', 'watch', 'templates', 'scripts', 'images', 'server']);
