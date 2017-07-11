import config from './rollup.config.js'

import buble from 'rollup-plugin-buble';
import uglify from 'rollup-plugin-uglify';
import {minify} from 'uglify-es';
import postcss from 'rollup-plugin-postcss';
import cssnano from 'cssnano';

config.plugins.push.apply(config.plugins, [
  buble({
    exclude: [ '**/*.css', '**/*.sss' ],
    transforms: { dangerousForOf: true }
  }),
 uglify({}, minify),
  postcss({
    extensions: ['.css', '.sss'],
    plugins: [cssnano()]
  })
]);

config.entry = 'dist/aot-prod-ie11-main.js'; // entry point for the application
config.dest = 'bundle/bundle-ie11.min.js';

export default config;
