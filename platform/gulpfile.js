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
var pkginfo = require('./package.json')

var paths = {
  public: 'public',
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

gulp.task('3rd', function() {
  var pkgs = [];
  for(var n in pkginfo.dependencies){
    pkgs.push('node_modules/'+n+"/**/*");
  }

   return gulp.src(pkgs, {"base":"."})
     .pipe(gulp.dest(paths.public));
});

gulp.task('scripts', function() {
  return gulp.src(paths.scripts)
    .pipe(sourcemaps.init())
    .pipe(coffee())
    .pipe(uglify())
    .pipe(concat('all.min.js'))
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

gulp.task('build', ['templates', 'scripts', '3rd', 'styles', 'images']);
