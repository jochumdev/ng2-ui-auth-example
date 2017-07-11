import './styles';
import './polyfills-ie11';
import './polyfills';
import {enableProdMode} from '@angular/core';
import {ClientModuleNgFactory} from './aot/src/client.module.ngfactory';
import {platformBrowser} from '@angular/platform-browser';

enableProdMode();

platformBrowser().bootstrapModuleFactory(ClientModuleNgFactory)
    .catch(err => console.error(err));
