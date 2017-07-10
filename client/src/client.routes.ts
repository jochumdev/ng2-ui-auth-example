import {ModuleWithProviders} from '@angular/core';
import {RouterModule} from '@angular/router';
import {AuthGuard} from './services/auth.guard';
import {ProfileComponent} from './components/profile.component';
import {LoginComponent} from './components/login.component';
import {SignupComponent} from './components/signup.component';
import {HomeComponent} from './components/home.component';
/**
 * Created by Ron on 03/10/2016.
 */
export const CLIENT_ROUTER_PROVIDERS = [
    AuthGuard
];
export const routing: ModuleWithProviders = RouterModule.forRoot([
    {
        path: '',
        redirectTo: 'home',
        pathMatch: 'full'
    },
    {
        path: 'profile',
        component: ProfileComponent,
        canActivate: [AuthGuard],
    },
    {
        path: 'login',
        component: LoginComponent,
    },
    {
        path: 'signup',
        component: SignupComponent,
    },
    {
        path: 'home',
        component: HomeComponent,
    }
]);
