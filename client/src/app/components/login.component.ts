import { Router } from '@angular/router';
import { Response } from '@angular/http';
import { AuthService } from 'ng2-ui-auth';
import { ILoginData } from '../interfaces';
import { Component, OnInit } from '@angular/core';
import { ErrorHandleService } from '../services/error-handle.service';
import { FormBuilder, FormControl, Validators, FormGroup } from '@angular/forms';
import { FormHelperService } from '../services/form-helper.service';
import { SettingsService } from '../services/settings.service';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/observable/throw';
/**
 * Created by Ron on 03/10/2016.
 */

@Component({
  selector: 'my-login',
  templateUrl: 'login.component.html',
})
export class LoginComponent implements OnInit {
  form: FormGroup;

  constructor(private auth: AuthService,
    private router: Router,
    private eh: ErrorHandleService,
    private fb: FormBuilder,
    private ss: SettingsService,
    public fh: FormHelperService) {
  }

  ngOnInit() {
    this.form = this.fb.group({
      username: new FormControl('', [Validators.required, Validators.minLength(3)]),
      password: new FormControl('', [Validators.required, Validators.minLength(6)]),
    })
  }

  login(loginData: ILoginData) {
    this.auth.login(loginData)
      .map((resp: Response) => resp)
      .catch((resp: Response) => {
        return Observable.throw(resp);
      })
      .subscribe({
        complete: () => {
          this.eh.saveMessage('success', 'You have successfully signed in!');
          this.router.navigateByUrl('profile');
        },
        error: (err: any) => this.eh.handleError(err),
      });
  }

  hasProvider(provider: string): boolean {
    return this.ss.hasProvider(provider);
  }

  externalLogin(provider: string) {
    this.auth.authenticate(provider)
      .map((resp: Response) => resp)
      .catch((resp: Response) => {
        return Observable.throw(resp);
      })
      .subscribe({
        complete: () => {
          this.eh.saveMessage('success', 'You have successfully signed in with ' + provider + '!');
          this.router.navigateByUrl('profile');
        },
        error: (err: any) => this.eh.handleError(err),
      });
  }
}
