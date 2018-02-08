import {AfterViewChecked, Component, OnInit, ViewContainerRef} from '@angular/core';
import {Response} from '@angular/http';
import {AuthService} from 'ng2-ui-auth';
import {Router} from '@angular/router';
import {ErrorHandleService} from './services/error-handle.service';
import {SettingsService, SettingAuthAllowSignup} from './services/settings.service';
import {ToastsManager} from 'ng2-toastr';
/**
 * Created by Ron on 03/10/2016.
 */
@Component({
    selector: 'app-root',
    templateUrl: 'app.component.html',
})
export class AppComponent implements OnInit, AfterViewChecked {
  constructor(private _vcr: ViewContainerRef,
              private _auth: AuthService,
              private _router: Router,
              private _eh: ErrorHandleService,
              private _ss: SettingsService,
              private _toastr: ToastsManager) {
    this._eh.setRootViewContainerRef(this._vcr);
  }

  ngOnInit() {
    this._ss.fetchSettings();

    // Allow toast on click dismiss, see:
    // https://github.com/PointInside/ng2-toastr/issues/61
    this._toastr.onClickToast()
      .subscribe( toast => {
        if (toast.timeoutId) {
          clearTimeout(toast.timeoutId);
        }
        this._toastr.dismissToast(toast)
      });
  }

  ngAfterViewChecked() {
    // You can call everywhere ErrorHandleService.saveMessage and .saveError
    // the code here will show them after the next routing happened.
    // This is usefull when you want to show a toastr message/error AFTER a routing.
    this._eh.showSavedMessage();
    this._eh.showSavedError();
  }

  isAuthenticated() {
    return this._auth.isAuthenticated()
  }

  signupAllowed() {
    if (this._ss.settings.has(SettingAuthAllowSignup)) {
      return this._ss.settings.get(SettingAuthAllowSignup).asBool() == true;
    }

    return false;
  }

  logout() {
    this._auth.logout()
        .subscribe({
            error: (err: any) => this._eh.saveError(<Response>err.json()),
            complete: () => {
              this._eh.saveMessage('info', 'You have been logged out');
              this._router.navigateByUrl('home');
            }
        });
  }
}
