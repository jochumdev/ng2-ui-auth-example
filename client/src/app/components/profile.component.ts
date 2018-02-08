import { AuthService } from 'ng2-ui-auth';
import { HttpClient } from '@angular/common/http';
import { OnInit, Component } from '@angular/core';
import { IProfileUser } from '../interfaces';
import { Observable } from 'rxjs/Observable';
import { FormBuilder, FormControl, Validators, FormGroup } from '@angular/forms';
import 'rxjs/add/operator/map';
import { ErrorHandleService } from '../services/error-handle.service';
import { FormHelperService } from '../services/form-helper.service';
import { SettingsService } from '../services/settings.service';
import { IProfileData } from '../interfaces';

import 'rxjs/add/observable/throw';

/**
 * Created by Ron on 03/10/2016.
 */

@Component({
  selector: 'my-profile',
  templateUrl: 'profile.component.html',
})
export class ProfileComponent implements OnInit {
  user: IProfileUser;
  form: FormGroup;

  constructor(private http: HttpClient,
    private auth: AuthService,
    private fb: FormBuilder,
    private eh: ErrorHandleService,
    private ss: SettingsService,
    public fh: FormHelperService) {
  }

  ngOnInit() {
    this.user = <IProfileUser>{};

    this.form = this.fb.group({
      displayName: new FormControl(),
      email: new FormControl('', [Validators.email]),
    });

    // Fetch the users profile
    this._updateUserData();
  }

  protected _updateUserData() {
    this.http.get('/api/auth/v1/me')
      .map((response: any) => response as IProfileUser)
      .catch((err: any) => {
        this.eh.handleError(err);
        return Observable.throw('Server error');
      })
      .subscribe((user: IProfileUser) => {
        this.user = user
      });
  }

  updateProfile(pdata: IProfileData) {
    this.http.put('/api/auth/v1/me', pdata)
      .map((response: any) => response as IProfileUser)
      .catch((err: any) => {
        this.eh.handleError(err);
        return Observable.throw('Server error');
      })
      .subscribe((user: IProfileUser) => {
        this.user = user;
        this.eh.saveMessage('success', 'Your profile has been updated');
      });
  }

  unlink(provider: string) {
    this.auth.unlink(provider)
      .map((response: any) => response)
      .subscribe({
        error: (err: any) => this.eh.handleError(err),
        complete: () => this._updateUserData()
      })
  }

  link(provider: string) {
    this.auth.link(provider)
      .map((response: any) => response)
      .subscribe({
        error: (err: any) => this.eh.handleError(err),
        complete: () => this._updateUserData()
      });
  }

  hasProvider(provider: string): boolean {
    return this.ss.hasProvider(provider);
  }
}
