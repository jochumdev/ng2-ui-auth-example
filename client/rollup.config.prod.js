import config from './rollup.config.js'

import buble from 'rollup-plugin-buble';
import uglify from 'rollup-plugin-uglify';
import {minify} from 'uglify-es';
import postcss from 'rollup-plugin-postcss';
import cssnano from 'cssnano';

config.plugins.push.apply(config.plugins, [
  buble({
    exclude: [ '**/*.css' ],
    transforms: { dangerousForOf: true }
  }),
  uglify({}, minify),
  postcss({
    extensions: ['.css', '.sss'],
    plugins: [cssnano()]
  })
]);

config.entry = 'dist/aot-prod-main.js'; // entry point for the application
config.dest = 'bundle/bundle.min.js';

export default config;
