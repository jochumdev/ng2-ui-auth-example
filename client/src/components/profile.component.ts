import {JwtHttp, AuthService} from 'ng2-ui-auth';
import {OnInit, Component} from '@angular/core';
import {Response} from '@angular/http';
import {IProfileUser} from '../interfaces';
import {Observable} from 'rxjs/Observable';
import {FormBuilder, FormControl, Validators, FormGroup} from '@angular/forms';
import 'rxjs/add/operator/map';
import {ErrorHandleService} from '../services/error-handle.service';
import {FormHelperService} from '../services/form-helper.service';
import {SettingsService} from '../services/settings.service';
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

    constructor(private http: JwtHttp,
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
        .map((response: Response) => <IProfileUser>response.json())
        .catch((err: any) => {
          this.eh.handleError(err);
          return Observable.throw('Server error');
        })
        .subscribe((user: IProfileUser) => {
            this.user = user
        });
    }

    unlink(provider: string) {
      this.auth.unlink(provider, {})
        .map((response: Response) => response.json())
        .subscribe({
            error: (err: any) => this.eh.handleError(err),
            complete: () => this._updateUserData()
        })
    }

    link(provider: string) {
      this.auth.link(provider)
        .map((response: Response) => response.json())
        .subscribe({
            error: (err: any) => this.eh.handleError(err),
            complete: () => this._updateUserData()
        });
    }

    hasProvider(provider: string): boolean {
      return this.ss.hasProvider(provider);
    }
}
