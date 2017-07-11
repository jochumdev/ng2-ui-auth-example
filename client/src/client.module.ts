import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {NgModule} from '@angular/core';
import {Ng2UiAuthModule} from 'ng2-ui-auth';
import {HttpModule, JsonpModule} from '@angular/http';
import {MyAuthConfig} from './config';
import {routing, CLIENT_ROUTER_PROVIDERS} from './client.routes';
import {ClientComponent} from './client.component';
import {ProfileComponent} from './components/profile.component';
import {LoginComponent} from './components/login.component';
import {SignupComponent} from './components/signup.component';
import {HomeComponent} from './components/home.component';
import {ReactiveFormsModule, FormsModule} from '@angular/forms';
import {ToastModule} from 'ng2-toastr';
import {ErrorHandleService} from './services/error-handle.service';
import {FormHelperService} from './services/form-helper.service';
import {SettingsService} from './services/settings.service';
/**
 * Created by Ron on 03/10/2016.
 */

@NgModule({
    imports: [
        BrowserModule,
        HttpModule,
        JsonpModule,
        FormsModule,
        ReactiveFormsModule,
        routing,
        Ng2UiAuthModule.forRoot(MyAuthConfig),
        ToastModule.forRoot(),
        BrowserAnimationsModule,
    ],
    providers: [
        SettingsService,
        ErrorHandleService,
        FormHelperService,
        CLIENT_ROUTER_PROVIDERS
    ],
    declarations: [
        ClientComponent,
        ProfileComponent,
        LoginComponent,
        SignupComponent,
        HomeComponent
    ],
    bootstrap: [
        ClientComponent
    ]
})
export class ClientModule {
}
