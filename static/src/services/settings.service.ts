import {Injectable} from '@angular/core';
import {Response} from '@angular/http';
import {ConfigService, JwtHttp} from 'ng2-ui-auth';
import {Observable} from 'rxjs/Observable';
import {ErrorHandleService} from './error-handle.service';
import {ISettingsProperty} from '../interfaces';

/**
 * Created by pcdummy on 07/03/2017.
 */
export const SettingAuthAllowSignup = "app_auth.AllowSignup";
export const SettingAuthGoogleClientID = "app_auth.GoogleClientID";
export const SettingAuthFacebookClientID = "app_auth.FacebookClientID";
export const SettingAuthGithubClientID = "app_auth.GithubClientID";

export class SettingsProperty implements ISettingsProperty {
  key: string;
  type: string;
  value: any;
  edit: boolean;

  constructor(key: string, type: string, value: any, edit: boolean) {
    this.key = key;
    this.type = type;
    this.value = value;
    this.edit = edit;
  }

  asString(): string {
    return <string>this.value;
  }

  asBool(): boolean {
    return <boolean>this.value;
  }
}

@Injectable()
export class SettingsService {
  settings: Map<string, SettingsProperty>;

  constructor(private _http: JwtHttp,
              private _eh: ErrorHandleService,
              private _ac: ConfigService) {
    this.settings = new Map<string, SettingsProperty>();
  }

  fetchSettings() {
    this._http.get('/api/settings/v1/')
      .map((response: Response) => <Array<SettingsProperty>>response.json())
      .catch((err: any) => {
        this._eh.handleError(err);
        return Observable.throw('Server error');
      })
      .subscribe((settings: Array<SettingsProperty>) => {
        settings.forEach((setting: SettingsProperty) => {
          this.settings.set(setting.key, new SettingsProperty(
            setting.key, setting.type, setting.value, setting.edit
          ));
        });

        if (this.settings.has(SettingAuthGoogleClientID)) {
          this._ac.providers['google'].clientId = this.settings.get(SettingAuthGoogleClientID).asString();
        }
        if (this.settings.has(SettingAuthFacebookClientID)) {
          this._ac.providers['facebook'].clientId = this.settings.get(SettingAuthFacebookClientID).asString();
        }
        if (this.settings.has(SettingAuthGithubClientID)) {
          this._ac.providers['github'].clientId = this.settings.get(SettingAuthGithubClientID).asString();
        }
      });
  }

  hasProvider(provider: string): boolean {
    if (provider == 'google') {
      if (this.settings.has(SettingAuthGoogleClientID)) {
        if (this.settings.get(SettingAuthGoogleClientID).asString() != "") {
          return true;
        }
      }
    } else if (provider == "facebook") {
      if (this.settings.has(SettingAuthFacebookClientID)) {
        if (this.settings.get(SettingAuthFacebookClientID).asString() != "") {
          return true;
        }
      }
    } else if (provider == "github") {
      if (this.settings.has(SettingAuthGithubClientID)) {
        if (this.settings.get(SettingAuthGithubClientID).asString() != "") {
          return true;
        }
      }
    }

    return false;
  }
}
