import './styles';
import './polyfills';
import {ClientModuleNgFactory} from './aot/src/client.module.ngfactory';
import {platformBrowser} from '@angular/platform-browser';

platformBrowser().bootstrapModuleFactory(ClientModuleNgFactory)
    .catch(err => console.error(err));
