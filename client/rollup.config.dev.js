import config from './rollup.config.js'

import buble from 'rollup-plugin-buble';
import sourcemaps from 'rollup-plugin-sourcemaps';
import postcss from 'rollup-plugin-postcss';

config.plugins.push.apply(config.plugins, [
  sourcemaps(),
  postcss({
    sourceMap: true,
    extensions: ['.css', '.sss']
  })
]);

config.entry = 'dist/aot-dev-main.js'; // entry point for the application
config.dest = 'bundle/bundle.js';

export default config;
