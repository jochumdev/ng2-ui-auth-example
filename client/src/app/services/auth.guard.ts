import { RouterStateSnapshot, ActivatedRouteSnapshot, Router, CanActivate } from '@angular/router';
import { AuthService } from 'ng2-ui-auth';
import { Injectable } from '@angular/core';

/**
 * Created by Ron on 03/10/2016.
 */
@Injectable()
export class AuthGuard implements CanActivate {
    constructor(private auth: AuthService, private router: Router) { }
    canActivate(
        _next: ActivatedRouteSnapshot,
        _state: RouterStateSnapshot
    ) {
        if (this.auth.isAuthenticated()) { return true; }
        this.router.navigateByUrl('login');
        return false;
    }
}
