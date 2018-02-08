import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { Ng2UiAuthModule } from 'ng2-ui-auth';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { myAuthConfig } from './config';
import { routing, CLIENT_ROUTER_PROVIDERS } from './app.routes';
import { AppComponent } from './app.component';
import { ProfileComponent } from './components/profile.component';
import { LoginComponent } from './components/login.component';
import { SignupComponent } from './components/signup.component';
import { HomeComponent } from './components/home.component';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { ToastModule } from 'ng2-toastr';
import { ErrorHandleService } from './services/error-handle.service';
import { FormHelperService } from './services/form-helper.service';
import { SettingsService } from './services/settings.service';
import { JsonInterceptorProvider } from './services/json-interceptor.service';
/**
 * Created by Ron on 03/10/2016.
 */


@NgModule({
    imports: [
        BrowserModule,
        // import HttpClientModule after BrowserModule.
        HttpClientModule,
        FormsModule,
        ReactiveFormsModule,
        routing,
        Ng2UiAuthModule.forRoot(myAuthConfig),
        ToastModule.forRoot(),
        BrowserAnimationsModule,
    ],
    providers: [
        JsonInterceptorProvider,
        SettingsService,
        ErrorHandleService,
        FormHelperService,
        CLIENT_ROUTER_PROVIDERS
    ],
    declarations: [
        AppComponent,
        ProfileComponent,
        LoginComponent,
        SignupComponent,
        HomeComponent
    ],
    bootstrap: [
        AppComponent
    ]
})
export class AppModule {
}
